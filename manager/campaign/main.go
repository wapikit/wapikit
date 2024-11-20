package campaign_manager

import (
	"database/sql"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	wapi "github.com/sarthakjdev/wapi.go/pkg/client"
	wapiComponents "github.com/sarthakjdev/wapi.go/pkg/components"
	"github.com/sarthakjdev/wapikit/internal/core/utils"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/sarthakjdev/wapikit/.db-generated/model"
	table "github.com/sarthakjdev/wapikit/.db-generated/table"
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

	// ! TODO: write a query which creates a sorted list of contacts aggregated from all the campaign list

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
}

func NewCampaignManager() *CampaignManager {
	return &CampaignManager{
		runningCampaigns:      make(map[string]*runningCampaign),
		runningCampaignsMutex: sync.RWMutex{},
		// ! TODO: set the message rate here, may be by fetching it from whatsapp api to get the limit allowed to the account in use
		messageQueue: make(chan *CampaignMessage),
		// 1000 campaigns can be queued at a time
		campaignQueue: make(chan *runningCampaign, 1000),
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
			WebhookPath:       "",
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

			campaignsQuery := SELECT(table.Campaign.AllColumns, table.WhatsappBusinessAccount.AllColumns).
				FROM(table.Campaign.LEFT_JOIN(
					table.WhatsappBusinessAccount, table.WhatsappBusinessAccount.OrganizationId.EQ(table.Campaign.OrganizationId),
				)).
				WHERE(table.Campaign.Status.EQ(utils.EnumExpression(model.CampaignStatus_Running.String())).
					AND(table.Campaign.UniqueId.NOT_IN(runningCampaignExpression...)))

			err := campaignsQuery.Query(cm.Db, &runningCampaigns)

			if err != nil {
				cm.Logger.Error("error fetching running campaigns from the database", err)
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

	// ! TODO: add the components to the template message

	messagingClient := client.NewMessagingClient(
		message.campaign.phoneNumberToUse,
	)
	_, err = messagingClient.Message.Send(templateMessage, message.contact.PhoneNumber)

	if err != nil {
		return err
	} else {
		message.campaign.lastContactIdSent = message.contact.UniqueId.String()
	}

	message.campaign.sent.Add(1)
	message.campaign.wg.Done() // * decrement the wg, because the message has been sent

	// ! TODO: update the database campaign with the last contact id sent the message to
	return nil
}

// this function gets called from the API server handlers, when user either pauses or cancels the campaign
func (cm *CampaignManager) StopCampaign(campaignUniqueId string) {
	cm.runningCampaignsMutex.RLock()
	if campaign, ok := cm.runningCampaigns[campaignUniqueId]; ok {
		campaign.stopCampaign()
	}
	cm.runningCampaignsMutex.RUnlock()
}
