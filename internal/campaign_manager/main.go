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

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/wapikit/wapikit/.db-generated/model"
	table "github.com/wapikit/wapikit/.db-generated/table"
	"github.com/wapikit/wapikit/internal/utils"
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
	WapiClient       *wapi.Client `json:"wapiclient"`
	PhoneNumberToUse string       `json:"phoneNumberToUse"`

	LastContactIdSent string       `json:"lastContactIdSent"`
	Sent              atomic.Int64 `json:"sent"`
	ErrorCount        atomic.Int64 `json:"errorCount"`

	IsStopped *atomic.Bool     `json:"isStopped"`
	Manager   *CampaignManager `json:"manager"`

	wg *sync.WaitGroup
}

// this function returns if the messages are exhausted or not
// if yes, then it will return false, and the campaign will be removed from the running campaigns list
func (rc *runningCampaign) nextContactsBatch() bool {
	var contacts []model.Contact

	contactsCte := CTE("contacts")
	updateCampaignLastContactSentIdCte := CTE("updateCampaignLastContactSentId")

	if rc.LastContactIdSent == "" {
		// assign a empty uuid here, so that the query can fetch the first contact
		rc.LastContactIdSent = uuid.MustParse("00000000-0000-0000-0000-000000000000").String()
	}

	lastContactSentUuid, err := uuid.Parse(rc.LastContactIdSent)

	if err != nil {
		rc.Manager.Logger.Error("error parsing lastContactSentUuid", err.Error())
		return false
	}

	campaignUniqueId, err := uuid.Parse(rc.UniqueId.String())

	if err != nil {
		rc.Manager.Logger.Error("error parsing campaignUniqueId", err.Error())
		return false
	}

	var contactLists []model.ContactList

	listIdsQuery := SELECT(table.ContactList.AllColumns, table.CampaignList.AllColumns).
		FROM(table.ContactList.INNER_JOIN(table.CampaignList, table.ContactList.UniqueId.EQ(table.CampaignList.ContactListId))).
		WHERE(table.CampaignList.CampaignId.EQ(UUID(campaignUniqueId)))

	err = listIdsQuery.Query(rc.Manager.Db, &contactLists)

	if err != nil {
		rc.Manager.Logger.Error("error fetching contact lists from the database", err.Error())
		return false
	}

	contactListIdExpression := make([]Expression, 0, len(contactLists))
	for _, contactList := range contactLists {
		contactListUuid, err := uuid.Parse(contactList.UniqueId.String())
		if err != nil {
			continue
		}
		contactListIdExpression = append(contactListIdExpression, UUID(contactListUuid))
	}

	var fromClause ReadableTable

	if len(contactListIdExpression) > 0 {
		fromClause = table.Contact.
			INNER_JOIN(
				table.ContactListContact, table.ContactListContact.ContactId.EQ(table.Contact.UniqueId).
					AND(table.ContactListContact.ContactListId.IN(contactListIdExpression...)),
			)
	} else {
		fromClause = table.Contact.
			INNER_JOIN(
				table.ContactListContact, table.ContactListContact.ContactId.EQ(table.Contact.UniqueId),
			)
	}

	nextContactsQuery := WITH(
		contactsCte.AS(
			SELECT(table.Contact.AllColumns, table.ContactListContact.AllColumns).
				FROM(fromClause).
				WHERE(table.Contact.UniqueId.GT(UUID(lastContactSentUuid))).
				DISTINCT(table.Contact.UniqueId).
				ORDER_BY(table.Contact.UniqueId).
				LIMIT(100),
		),
		updateCampaignLastContactSentIdCte.AS(
			table.Campaign.UPDATE(table.Campaign.LastContactSent).
				WHERE(table.Campaign.UniqueId.EQ(UUID(campaignUniqueId))).
				SET(UUID(lastContactSentUuid)),
		),
	)(
		SELECT(
			contactsCte.AllColumns(),
		).FROM(
			contactsCte,
		),
	)

	err = nextContactsQuery.Query(rc.Manager.Db, &contacts)

	if err != nil {
		rc.Manager.Logger.Error("error fetching contacts from the database", err.Error(), nil)
		return false
	}

	// * all contacts have been sent the message, so return false
	if len(contacts) == 0 {
		return false
	}

	for _, contact := range contacts {
		// * add the message to the message queue
		message := &CampaignMessage{
			Campaign: rc,
			Contact:  contact,
		}

		select {
		case rc.Manager.messageQueue <- message:
			rc.wg.Add(1)
		default:
			// * if the message queue is full, then return true, so that the campaign can be queued again
			return true
		}
	}

	return false
}

func (rc *runningCampaign) stop() {
	if rc.IsStopped.Load() {
		return
	}
	rc.IsStopped.Store(true)
}

// this function will only run when the campaign is exhausted its subscriber list
func (rc *runningCampaign) cleanUp() {
	defer func() {
		rc.Manager.runningCampaignsMutex.Lock()
		delete(rc.Manager.runningCampaigns, rc.UniqueId.String())
		rc.Manager.runningCampaignsMutex.Unlock()
	}()

	// check the fresh status of the campaign, if it is still running, then update the status to finished
	var campaign model.Campaign

	campaignQuery := SELECT(table.Campaign.AllColumns).
		FROM(table.Campaign).
		WHERE(table.Campaign.UniqueId.EQ(String(rc.UniqueId.String())))

	err := campaignQuery.Query(rc.Manager.Db, &campaign)

	if err != nil {
		rc.Manager.Logger.Error("error fetching campaign from the database", err.Error(), nil)
		// campaign not found in the db for some reason, it will be removed from the running campaigns list
		return
	}

	if campaign.Status == model.CampaignStatusEnum_Running {
		_, err = rc.Manager.updatedCampaignStatus(rc.UniqueId.String(), model.CampaignStatusEnum_Finished)
		if err != nil {
			rc.Manager.Logger.Error("error updating campaign status", err.Error(), nil)
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
	Campaign *runningCampaign `json:"campaign"`
	Contact  model.Contact    `json:"contact"`
}

func (cm *CampaignManager) newRunningCampaign(dbCampaign model.Campaign, businessAccount model.WhatsappBusinessAccount) *runningCampaign {
	campaign := runningCampaign{
		Campaign: dbCampaign,
		WapiClient: wapi.New(&wapi.ClientConfig{
			BusinessAccountId: businessAccount.AccountId,
			ApiAccessToken:    businessAccount.AccessToken,
			WebhookSecret:     businessAccount.WebhookSecret,
		}),
		PhoneNumberToUse:  dbCampaign.PhoneNumber,
		LastContactIdSent: "",
		Sent:              atomic.Int64{},
		ErrorCount:        atomic.Int64{},
		Manager:           cm,
		wg:                &sync.WaitGroup{},
		IsStopped:         &atomic.Bool{},
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

	cm.Logger.Info("campaign manager started.")

	// * process the campaign queue, means listen to the campaign queue, and then for each campaign, call the function to next subscribers
	for campaign := range cm.campaignQueue {
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

func (cm *CampaignManager) updatedCampaignStatus(campaignId string, status model.CampaignStatusEnum) (bool, error) {
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
			for _, campaignId := range currentRunningCampaignIds {
				campaignUuid, err := uuid.Parse(campaignId)
				if err != nil {
					cm.Logger.Error("error parsing campaign id", err.Error())
					continue
				}
				runningCampaignExpression = append(runningCampaignExpression, UUID(campaignUuid))
			}

			whereCondition := table.Campaign.Status.EQ(utils.EnumExpression(model.CampaignStatusEnum_Running.String()))

			if len(runningCampaignExpression) > 0 {
				whereCondition = table.Campaign.Status.EQ(utils.EnumExpression(model.CampaignStatusEnum_Running.String()))
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

			if message.Campaign.IsStopped.Load() {
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
			// // ! TODO: send an update to the websocket server, updating the count of messages sent for the campaign
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
	client := message.Campaign.WapiClient

	templateInUse, err := client.Business.Template.Fetch(*message.Campaign.MessageTemplateId)

	if err != nil {
		message.Campaign.ErrorCount.Add(1)
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
				if len(button.Example) > 0 {
					doTemplateRequireParameter = true
				}
			}
		}
	}

	type templateComponentParameters struct {
		Header  []string `json:"header"`
		Body    []string `json:"body"`
		Buttons []string `json:"buttons"`
	}

	var parameterStoredInDb templateComponentParameters
	err = json.Unmarshal([]byte(*message.Campaign.TemplateMessageComponentParameters), &parameterStoredInDb)
	if err != nil {
		return fmt.Errorf("error unmarshalling template parameters: %v", err)
	}

	// Check if the struct is at its zero value
	if doTemplateRequireParameter && reflect.DeepEqual(parameterStoredInDb, templateComponentParameters{}) {
		// Stop the campaign and return an error
		cm.StopCampaign(message.Campaign.UniqueId.String())
		return fmt.Errorf("template requires parameters, but no parameter found in the database")
	}

	for _, component := range templateInUse.Components {
		switch component.Type {
		case "BODY":
			{
				if len(component.Example.BodyText) > 0 {
					bodyParameters := []wapiComponents.TemplateMessageParameter{}
					for _, bodyText := range parameterStoredInDb.Body {
						bodyParameters = append(bodyParameters, wapiComponents.TemplateMessageBodyAndHeaderParameter{
							Type: wapiComponents.TemplateMessageParameterTypeText,
							Text: &bodyText,
						})
					}
					templateMessage.AddBody(wapiComponents.TemplateMessageComponentBodyType{
						Type:       wapiComponents.TemplateMessageComponentTypeBody,
						Parameters: bodyParameters,
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
					// use header text
					headerParameters := []wapiComponents.TemplateMessageParameter{}
					for _, headerText := range parameterStoredInDb.Header {
						headerParameters = append(headerParameters, wapiComponents.TemplateMessageBodyAndHeaderParameter{
							Type: wapiComponents.TemplateMessageParameterTypeText,
							Text: &headerText,
						})
					}
					templateMessage.AddHeader(wapiComponents.TemplateMessageComponentHeaderType{
						Type:       wapiComponents.TemplateMessageComponentTypeHeader,
						Parameters: headerParameters,
					})
				} else if component.Format == "IMAGE" {
					headerParameters := []wapiComponents.TemplateMessageParameter{}
					for _, mediaUrl := range parameterStoredInDb.Header {
						headerParameters = append(headerParameters, wapiComponents.TemplateMessageBodyAndHeaderParameter{
							Type: wapiComponents.TemplateMessageParameterTypeText,
							Image: &wapiComponents.TemplateMessageParameterMedia{
								Link: mediaUrl,
							},
						})
					}
					templateMessage.AddHeader(wapiComponents.TemplateMessageComponentHeaderType{
						Type:       wapiComponents.TemplateMessageComponentTypeHeader,
						Parameters: headerParameters,
					})

					// ! TODO: use header handle
					// MessageTemplateComponentFormatDocument MessageTemplateComponentFormat = "DOCUMENT"
					// MessageTemplateComponentFormatVideo    MessageTemplateComponentFormat = "VIDEO"
					// MessageTemplateComponentFormatLocation MessageTemplateComponentFormat = "LOCATION"

				} else if component.Format == "VIDEO" {
					headerParameters := []wapiComponents.TemplateMessageParameter{}
					for _, mediaUrl := range parameterStoredInDb.Header {
						headerParameters = append(headerParameters, wapiComponents.TemplateMessageBodyAndHeaderParameter{
							Type: wapiComponents.TemplateMessageParameterTypeVideo,
							Image: &wapiComponents.TemplateMessageParameterMedia{
								Link: mediaUrl,
							},
						})
					}
					templateMessage.AddHeader(wapiComponents.TemplateMessageComponentHeaderType{
						Type:       wapiComponents.TemplateMessageComponentTypeHeader,
						Parameters: headerParameters,
					})

				} else if component.Format == "DOCUMENT" {
					headerParameters := []wapiComponents.TemplateMessageParameter{}
					for _, mediaUrl := range parameterStoredInDb.Header {
						headerParameters = append(headerParameters, wapiComponents.TemplateMessageBodyAndHeaderParameter{
							Type: wapiComponents.TemplateMessageParameterTypeDocument,
							Document: &wapiComponents.TemplateMessageParameterMedia{
								Link: mediaUrl,
							},
						})
					}
					templateMessage.AddHeader(wapiComponents.TemplateMessageComponentHeaderType{
						Type:       wapiComponents.TemplateMessageComponentTypeHeader,
						Parameters: headerParameters,
					})
				} else if component.Format == "LOCATION" {

					// ! TODO: implement location type here

					// headerParameters := []wapiComponents.TemplateMessageParameter{}
					// for _, mediaUrl := range parameterStoredInDb.Header {
					// 	headerParameters = append(headerParameters, wapiComponents.TemplateMessageBodyAndHeaderParameter{
					// 		Type: wapiComponents.TemplateMessageParameterTypeLocation,
					// 		Document: &wapiComponents.TemplateMessageParameterLocation{
					// 			Latitude: "0.0",
					// 		},
					// 	})
					// }
					// templateMessage.AddHeader(wapiComponents.TemplateMessageComponentHeaderType{
					// 	Type:       wapiComponents.TemplateMessageComponentTypeHeader,
					// 	Parameters: headerParameters,
					// })
				}

			}

		case "BUTTONS":
			{
				for index, button := range component.Buttons {
					switch button.Type {
					case "URL":
						{
							templateMessage.AddButton(wapiComponents.TemplateMessageComponentButtonType{
								Type:    wapiComponents.TemplateMessageComponentTypeButton,
								SubType: wapiComponents.TemplateMessageButtonComponentTypeUrl,
								Index:   index,
								Parameters: []wapiComponents.TemplateMessageParameter{
									wapiComponents.TemplateMessageButtonParameter{
										Type: wapiComponents.TemplateMessageButtonParameterTypeText,
										Text: parameterStoredInDb.Buttons[index],
									},
								}})
						}
					case "QUICK_REPLY":
						{
							templateMessage.AddButton(wapiComponents.TemplateMessageComponentButtonType{
								Type:    wapiComponents.TemplateMessageComponentTypeButton,
								Index:   index,
								SubType: wapiComponents.TemplateMessageButtonComponentTypeQuickReply,
								Parameters: []wapiComponents.TemplateMessageParameter{
									wapiComponents.TemplateMessageButtonParameter{
										Type:    wapiComponents.TemplateMessageButtonParameterTypePayload,
										Payload: parameterStoredInDb.Buttons[index],
									},
								},
							})
						}
					case "COPY_CODE":
						{
							// ! TODO: implement copy code button here
						}
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
		message.Campaign.PhoneNumberToUse,
	)

	response, err := messagingClient.Message.Send(templateMessage, message.Contact.PhoneNumber)
	fmt.Println("response", response)

	if err != nil {
		fmt.Errorf("error sending message to user: %v", err.Error())
		message.Campaign.ErrorCount.Add(1)
		return err
	}

	jsonMessage, err := templateMessage.ToJson(wapiComponents.ApiCompatibleJsonConverterConfigs{
		SendToPhoneNumber: message.Contact.PhoneNumber,
	})

	stringifiedJsonMessage := string(jsonMessage)
	fmt.Println("stringifiedJsonMessage", stringifiedJsonMessage)

	if err != nil {
		return err
	} else {
		message.Campaign.LastContactIdSent = message.Contact.UniqueId.String()
		// create a record in the db
		messageSent := model.Message{
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			CampaignId:      &message.Campaign.UniqueId,
			Direction:       model.MessageDirectionEnum_OutBound,
			ContactId:       message.Contact.UniqueId,
			PhoneNumberUsed: message.Campaign.PhoneNumberToUse,
			OrganizationId:  message.Campaign.OrganizationId,
			MessageData:     &stringifiedJsonMessage,
			MessageType:     model.MessageTypeEnum_Text,
			Status:          model.MessageStatusEnum_Delivered,
		}

		messageSentRecordQuery := table.Message.
			INSERT().
			MODEL(messageSent).
			RETURNING(table.Message.AllColumns)

		err := messageSentRecordQuery.Query(cm.Db, &messageSent)

		if err != nil {
			cm.Logger.Error("error saving message record to the database", err.Error())
		}
	}

	message.Campaign.Manager.rateLimiter.Incr(1)
	message.Campaign.Sent.Add(1)
	message.Campaign.wg.Done() // * decrement the wg, because the message has been sent

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
