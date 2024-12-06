package campaign_manager

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/paulbellamy/ratecounter"
	wapi "github.com/wapikit/wapi.go/pkg/client"
	wapiComponents "github.com/wapikit/wapi.go/pkg/components"
	"github.com/wapikit/wapikit/internal/core/utils"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/wapikit/wapikit/.db-generated/model"
	table "github.com/wapikit/wapikit/.db-generated/table"
)

// ! NOTE:
// ! for now this campaign manager does not adheres the whatsapp business conversation initiate rate limit and template message rate limit, which are unique for each business account.
// ! soon, in the later version we will fix this, but as of now it is required by the campaign admin to watch out for their own limit.
// ! it adheres a global 80 message per second limit for the whatsapp business api, which is the default limit for the whatsapp business api
// ! https://developers.facebook.com/docs/whatsapp/cloud-api/overview/#rate-limits

// ! also, there is a pair rate limit on wsp biz API, https://developers.facebook.com/docs/whatsapp/cloud-api/overview/#pair-rate-limits
// ! it checks that only up to 6 messages can be sent to a whatsapp phone number in a second, and up to 600 messages in a 24 per hour

var (
	messagesPerSecondLimit = 80
)

type runningCampaign struct {
	model.Campaign
	wapiClient       *wapi.Client
	phoneNumberToUse string

	lastContactIdSent string
	sent              atomic.Int64
	errorCount        atomic.Int64

	isStopped *atomic.Bool
	manager   *CampaignManager

	wg *sync.WaitGroup
}

// this function returns if the messages are exhausted or not
// if yes, then it will return false, and the campaign will be removed from the running campaigns list
func (rc *runningCampaign) nextContactsBatch() bool {
	var contacts []model.Contact

	campaignListCte := CTE("campaignLists")
	contactsCte := CTE("contacts")
	updateCampaignLastContactSentIdCte := CTE("updateCampaignLastContactSentId")

	contactListIds := WITH(
		campaignListCte.AS(
			SELECT(table.ContactList.AllColumns, table.CampaignList.AllColumns).
				FROM(table.CampaignList.LEFT_JOIN(table.ContactList, table.ContactList.UniqueId.EQ(table.CampaignList.ContactListId))).
				WHERE(table.CampaignList.CampaignId.EQ(String(rc.UniqueId.String()))),
		),
		contactsCte.AS(
			SELECT(table.Contact.AllColumns).
				FROM(table.Contact.LEFT_JOIN(table.ContactListContact, table.ContactListContact.ContactId.EQ(table.Contact.UniqueId))).
				WHERE(
					table.ContactListContact.ContactListId.IN(campaignListCte.SELECT(table.ContactList.UniqueId)).
						AND(table.Contact.UniqueId.GT(String(rc.lastContactIdSent))),
				).DISTINCT(table.Contact.UniqueId).
				ORDER_BY(table.Contact.UniqueId).
				LIMIT(100),
		),
		updateCampaignLastContactSentIdCte.AS(
			table.Campaign.UPDATE(table.Campaign.LastContactSent).
				WHERE(table.Campaign.UniqueId.EQ(UUID(rc.UniqueId))).
				SET(UUID(rc.LastContactSent)),
		),
	)(
		SELECT(
			contactsCte.SELECT(table.Contact.AllColumns),
		).FROM(
			contactsCte,
		),
	)

	err := contactListIds.Query(rc.manager.Db, &contacts)

	if err != nil {
		rc.manager.Logger.Error("error fetching contacts from the database", err.Error())
		return false
	}

	// * all contacts have been sent the message, so return false
	if len(contacts) == 0 {
		return false
	}

	for _, contact := range contacts {
		// * add the message to the message queue
		message := &CampaignMessage{
			campaign: rc,
			contact:  contact,
		}

		select {
		case rc.manager.messageQueue <- message:
			rc.wg.Add(1)
		default:
			// * if the message queue is full, then return true, so that the campaign can be queued again
			return true
		}
	}

	return false
}

func (rc *runningCampaign) stop() {
	if rc.isStopped.Load() {
		return
	}
	rc.isStopped.Store(true)
}

// this function will only run when the campaign is exhausted its subscriber list
func (rc *runningCampaign) cleanUp() {
	defer func() {
		rc.manager.runningCampaignsMutex.Lock()
		delete(rc.manager.runningCampaigns, rc.UniqueId.String())
		rc.manager.runningCampaignsMutex.Unlock()
	}()

	// check the fresh status of the campaign, if it is still running, then update the status to finished
	var campaign model.Campaign

	campaignQuery := SELECT(table.Campaign.AllColumns).
		FROM(table.Campaign).
		WHERE(table.Campaign.UniqueId.EQ(String(rc.UniqueId.String())))

	err := campaignQuery.Query(rc.manager.Db, &campaign)

	if err != nil {
		rc.manager.Logger.Error("error fetching campaign from the database", err.Error())
		// campaign not found in the db for some reason, it will be removed from the running campaigns list
		return
	}

	if campaign.Status == model.CampaignStatus_Running {
		_, err = rc.manager.updatedCampaignStatus(rc.UniqueId.String(), model.CampaignStatus_Finished)
		if err != nil {
			rc.manager.Logger.Error("error updating campaign status", err.Error())
		}
	}
}

type CampaignManager struct {
	Db     *sql.DB
	Logger slog.Logger

	runningCampaigns      map[string]*runningCampaign
	runningCampaignsMutex sync.RWMutex

	messageQueue  chan *CampaignMessage
	campaignQueue chan *runningCampaign

	rateLimiter *ratecounter.RateCounter
}

func NewCampaignManager(db *sql.DB, logger slog.Logger) *CampaignManager {
	return &CampaignManager{
		Db:     db,
		Logger: logger,

		runningCampaigns:      make(map[string]*runningCampaign),
		runningCampaignsMutex: sync.RWMutex{},
		// ! TODO: set the message rate here, may be by fetching it from whatsapp api to get the limit allowed to the account in use
		messageQueue: make(chan *CampaignMessage),
		// 1000 campaigns can be queued at a time
		campaignQueue: make(chan *runningCampaign, 1000),

		rateLimiter: ratecounter.NewRateCounter(1 * time.Second),
	}
}

type CampaignMessage struct {
	campaign *runningCampaign
	contact  model.Contact
}

func (cm *CampaignManager) newRunningCampaign(dbCampaign model.Campaign, businessAccount model.WhatsappBusinessAccount) *runningCampaign {
	campaign := runningCampaign{
		Campaign: dbCampaign,
		wapiClient: wapi.New(&wapi.ClientConfig{
			BusinessAccountId: businessAccount.AccountId,
			ApiAccessToken:    businessAccount.AccessToken,
			WebhookSecret:     businessAccount.WebhookSecret,
		}),
		phoneNumberToUse:  dbCampaign.PhoneNumber,
		lastContactIdSent: "",
		sent:              atomic.Int64{},
		errorCount:        atomic.Int64{},
		manager:           cm,
		wg:                &sync.WaitGroup{},
		isStopped:         &atomic.Bool{},
	}

	// * add the campaign to the wait group, because we are having a asynchronous setup for processing the messages of the campaign
	campaign.wg.Add(1)

	go func() {
		campaign.wg.Wait()
		campaign.stop()
		campaign.cleanUp()
	}()

	cm.runningCampaignsMutex.Lock()
	cm.runningCampaigns[campaign.UniqueId.String()] = &campaign
	cm.runningCampaignsMutex.Unlock()
	return &campaign
}

// Run starts the campaign manager
// main blocking function must be executed in a go routine
func (cm *CampaignManager) Run() {
	// * scan for campaign status changes every 5 seconds
	go cm.scanCampaigns()

	// * this function will process the message queue
	go cm.processMessageQueue()

	// * process the campaign queue, means listen to the campaign queue, and then for each campaign, call the function to next subscribers
	for _, campaign := range cm.runningCampaigns {
		hasContactsRemainingInQueue := campaign.nextContactsBatch()
		if hasContactsRemainingInQueue {
			// queue it again
			select {
			case cm.campaignQueue <- campaign:
			default:
			}
		} else {
			campaign.wg.Done()
		}
	}
}

func (cm *CampaignManager) updatedCampaignStatus(campaignId string, status model.CampaignStatus) (bool, error) {
	campaignUpdateQuery := table.Campaign.UPDATE(table.Campaign.Status).
		SET(status).
		WHERE(table.Campaign.UniqueId.EQ(String(campaignId)))

	_, err := campaignUpdateQuery.Exec(cm.Db)

	if err != nil {
		cm.Logger.Error("error updating campaign status", err.Error())
		return false, err
	}

	return true, nil
}

func (cm *CampaignManager) scanCampaigns() {

	// * scan for campaign status changes every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			currentRunningCampaignIds := cm.getRunningCampaignsUniqueIds()
			var runningCampaigns []struct {
				model.Campaign
				model.WhatsappBusinessAccount
			}

			runningCampaignExpression := make([]Expression, 0, len(currentRunningCampaignIds))
			for i, campaignId := range currentRunningCampaignIds {
				campaignUuid, err := uuid.Parse(campaignId)
				if err != nil {
					cm.Logger.Error("error parsing campaign id", err.Error())
					continue
				}
				runningCampaignExpression[i] = UUID(campaignUuid)
			}

			whereCondition := table.Campaign.Status.EQ(utils.EnumExpression(model.CampaignStatus_Running.String()))

			if len(runningCampaignExpression) > 0 {
				whereCondition = table.Campaign.Status.EQ(utils.EnumExpression(model.CampaignStatus_Running.String()))
			}

			campaignsQuery := SELECT(table.Campaign.AllColumns, table.WhatsappBusinessAccount.AllColumns).
				FROM(table.Campaign.LEFT_JOIN(
					table.WhatsappBusinessAccount, table.WhatsappBusinessAccount.OrganizationId.EQ(table.Campaign.OrganizationId),
				)).
				WHERE(whereCondition)

			context := context.Background()
			err := campaignsQuery.QueryContext(context, cm.Db, &runningCampaigns)

			if err != nil {
				cm.Logger.Error("error fetching running campaigns from the database", err)
			}

			if len(runningCampaigns) == 0 {
				// no running campaign found
				continue
			}

			for _, campaign := range runningCampaigns {
				campaignToAdd := cm.newRunningCampaign(campaign.Campaign, campaign.WhatsappBusinessAccount)
				select {
				case cm.campaignQueue <- campaignToAdd:
				default:
				}
			}
		}
	}
}

func (cm *CampaignManager) processMessageQueue() {
	for {
		select {
		case message, ok := <-cm.messageQueue:
			if !ok {
				return
			}

			if message.campaign.isStopped.Load() {
				// campaign has been stopped, so skip this message
				continue
			}

			// Check the rate limiter
			if cm.rateLimiter.Rate() >= int64(messagesPerSecondLimit) {
				// Rate limit exceeded, requeue the message and wait
				// if in case the queue will be full the below select walk will wait until the queue gets empty to add the message to the queue asynchronously
				select {
				case cm.messageQueue <- message:
				default:
				}
				time.Sleep(10 * time.Millisecond) // Sleep for a short duration before retrying
				continue
			}

			err := cm.sendMessage(message)
			// ! TODO: send an update to the websocket server, updating the count of messages sent for the campaign
			if err != nil {
				cm.Logger.Error("error sending message to user", err.Error())
				// ! TODO: broadcast this message to websocket via the API server event
			}
		}
	}
}

func (cm *CampaignManager) getRunningCampaignsUniqueIds() []string {
	cm.runningCampaignsMutex.RLock()
	uniqueIds := make([]string, 0, len(cm.campaignQueue))
	for _, campaign := range cm.runningCampaigns {
		uniqueIds = append(uniqueIds, campaign.UniqueId.String())
	}
	cm.runningCampaignsMutex.RUnlock()
	return uniqueIds
}

func (cm *CampaignManager) sendMessage(message *CampaignMessage) error {
	client := message.campaign.wapiClient

	templateInUse, err := client.Business.Template.Fetch(*message.campaign.MessageTemplateId)

	if err != nil {
		message.campaign.errorCount.Add(1)
		return fmt.Errorf("error fetching template: %v", err)
	}

	// * create the template message
	templateMessage, err := wapiComponents.NewTemplateMessage(
		&wapiComponents.TemplateMessageConfigs{
			Name:     templateInUse.Name,
			Language: templateInUse.Language,
		},
	)

	if err != nil {
		return fmt.Errorf("error creating template message: %v", err)
	}

	// * check if this template required parameters, if yes then check if we have parameter in db, else ignore, and if no parameter in db, then error

	doTemplateRequireParameter := false

	for _, component := range templateInUse.Components {

		if len(component.Example.BodyText) > 0 || len(component.Example.HeaderText) > 0 || len(component.Example.HeaderText) > 0 {
			doTemplateRequireParameter = true
		}

		if len(component.Buttons) > 0 {
			for _, button := range component.Buttons {
				if button.Example != "" {
					doTemplateRequireParameter = true
				}
			}
		}
	}

	type templateComponentParameters struct {
		Header []string `json:"header"`
		Body   []string `json:"body"`
		Button []string `json:"button"`
	}

	var parameterStoredInDb templateComponentParameters

	err = json.Unmarshal([]byte(*message.campaign.TemplateMessageComponentParameters), &parameterStoredInDb)
	if err != nil {
		return fmt.Errorf("error unmarshalling template parameters: %v", err)
	}

	// Check if the struct is at its zero value
	if doTemplateRequireParameter && reflect.DeepEqual(parameterStoredInDb, templateComponentParameters{}) {
		// Stop the campaign and return an error
		cm.StopCampaign(message.campaign.UniqueId.String())
		return fmt.Errorf("template requires parameters, but no parameter found in the database")
	}

	// ! TODO: add the components to the template message
	for _, component := range templateInUse.Components {
		switch component.Type {
		case "BODY":
			{
				if len(component.Example.BodyText) > 0 {
					templateMessage.AddBody(wapiComponents.TemplateMessageComponentBodyType{
						Type:       wapiComponents.TemplateMessageComponentTypeBody,
						Parameters: []wapiComponents.TemplateMessageParameter{},
					})
				} else {
					templateMessage.AddBody(wapiComponents.TemplateMessageComponentBodyType{
						Type:       wapiComponents.TemplateMessageComponentTypeBody,
						Parameters: []wapiComponents.TemplateMessageParameter{},
					})
				}

			}

		case "HEADER":
			{

				if len(component.Example.HeaderText) == 0 && len(component.Example.HeaderHandle) == 0 {
					templateMessage.AddHeader(wapiComponents.TemplateMessageComponentHeaderType{
						Type: wapiComponents.TemplateMessageComponentTypeHeader,
					})
				}

				if component.Format == "TEXT" {
					// use  header text
					templateMessage.AddHeader(wapiComponents.TemplateMessageComponentHeaderType{
						Type: wapiComponents.TemplateMessageComponentTypeHeader,
						Parameters: []wapiComponents.TemplateMessageParameter{
							wapiComponents.TemplateMessageBodyAndHeaderParameter{
								Type: wapiComponents.TemplateMessageParameterTypeText,
								// Text: "", // ! TODO: get this from the campaign db record user provided string
							},
						},
					})

				} else {
					// ! TODO: use header handle

				}

				// if example has length, it means the header has parameters to add
				if len(component.Example.BodyText) > 0 || len(component.Example.HeaderText) > 0 || len(component.Example.BodyText) > 0 {
					templateMessage.AddHeader(wapiComponents.TemplateMessageComponentHeaderType{
						Type: wapiComponents.TemplateMessageComponentTypeHeader,
					})
				} else {

				}

			}

		case "BUTTONS":
			{
				for _, button := range component.Buttons {
					switch button.Type {
					case "URL":
						{
							templateMessage.AddButton(wapiComponents.TemplateMessageComponentButtonType{
								Type:    wapiComponents.TemplateMessageComponentTypeButton,
								SubType: wapiComponents.TemplateMessageButtonComponentTypeUrl,
							})
						}
					case "QUICK_REPLY":
						{
							templateMessage.AddButton(wapiComponents.TemplateMessageComponentButtonType{
								Type:    wapiComponents.TemplateMessageComponentTypeButton,
								SubType: wapiComponents.TemplateMessageButtonComponentTypeQuickReply,
							})
						}
						// case "PHONE_NUMBER":
						// 	{
						// 		templateMessage.AddButton(wapiComponents.TemplateMessageComponentButtonType{
						// 			Type:    wapiComponents.TemplateMessageComponentTypeButton,
						// 			SubType: wapiComponents.TemplateMessageButtonComponentTyp,
						// 		})
						// 	}
					}

				}
			}

		case "FOOTER":
			{
				// ! TODO: to be implemented
			}

		}
	}

	messagingClient := client.NewMessagingClient(
		message.campaign.phoneNumberToUse,
	)
	_, err = messagingClient.Message.Send(templateMessage, message.contact.PhoneNumber)

	jsonMessage, err := templateMessage.ToJson(wapiComponents.ApiCompatibleJsonConverterConfigs{
		SendToPhoneNumber: "",
	})

	stringifiedJsonMessage := string(jsonMessage)

	if err != nil {
		return err
	} else {
		message.campaign.lastContactIdSent = message.contact.UniqueId.String()

		// create a record in the db
		messageSent := model.Message{
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			CampaignId:      &message.campaign.UniqueId,
			Direction:       model.MessageDirection_OutBound,
			ContactId:       message.contact.UniqueId,
			PhoneNumberUsed: message.campaign.phoneNumberToUse,
			OrganizationId:  message.campaign.OrganizationId,
			Content:         &stringifiedJsonMessage,
			Status:          model.MessageStatus_Delivered,
		}

		messageSentRecordQuery := table.Message.INSERT().
			MODEL(messageSent).
			RETURNING(table.Message.AllColumns)

		err := messageSentRecordQuery.Query(cm.Db, &messageSent)

		if err != nil {
			cm.Logger.Error("error saving message record to the database", err.Error())
		}
	}

	message.campaign.manager.rateLimiter.Incr(1)
	message.campaign.sent.Add(1)
	message.campaign.wg.Done() // * decrement the wg, because the message has been sent

	// ! TODO: update the database campaign with the last contact id sent the message to
	return nil
}

// this function gets called from the API server handlers, when user either pauses or cancels the campaign
func (cm *CampaignManager) StopCampaign(campaignUniqueId string) {
	cm.runningCampaignsMutex.RLock()
	if campaign, ok := cm.runningCampaigns[campaignUniqueId]; ok {
		campaign.stop()
	}
	cm.runningCampaignsMutex.RUnlock()
}
