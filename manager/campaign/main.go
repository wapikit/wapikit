package manager

import (
	"time"

	wapi "github.com/sarthakjdev/wapi.go/pkg/client"
	"github.com/sarthakjdev/wapikit/internal/core/utils"
	"github.com/sarthakjdev/wapikit/internal/interfaces"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/sarthakjdev/wapikit/.db-generated/model"
	table "github.com/sarthakjdev/wapikit/.db-generated/table"
)

// - it checks for the running status campaigns in db
// - it queues the message to be sent to user for the campaigns that are running
// - it checks if the campaign has ended and updates the status in db
// - on every message send it updates the last user id for the campaign in the db
// - it fetches the next batch of users to send the message to
// - it updates the campaign status to completed in db
// - it runs all this in memory
// - it must be executed in a go routine because it a long running blocking function, which continuously check for campaign and messages to be sent.

type CampaignManager struct {
	app interfaces.App

	messageQueue chan CampaignMessage

	activeCampaigns map[string]struct {
		campaign   model.Campaign
		wapiClient *wapi.Client
	}
}

type CampaignMessage struct {
	phoneNumberToUse string
	messageJson      string
	wapiClient       *wapi.Client
	contact          model.Contact
}

// each campaign will have its own wapi client with pre-feeded api access token and phone number

func NewCampaignManager() *CampaignManager {
	return &CampaignManager{}
}

// Run starts the campaign manager
// main blocking function must be executed in a go routine
func (cm *CampaignManager) Run() {
	// * scan for campaign status changes every 5 seconds
	go cm.scanCampaigns()

	// * this function will process the message queue
	go cm.processMessageQueue()
}

func (cm *CampaignManager) scanCampaigns() {
	// * scan for campaign status changes every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			var runningCampaigns []struct {
				model.Campaign
				model.WhatsappBusinessAccount
			}

			campaignsQuery := SELECT(table.Campaign.AllColumns, table.WhatsappBusinessAccount.AllColumns).
				FROM(table.Campaign.LEFT_JOIN(
					table.WhatsappBusinessAccount, table.WhatsappBusinessAccount.OrganizationId.EQ(table.Campaign.OrganizationId),
				)).
				WHERE(table.Campaign.Status.EQ(utils.EnumExpression(model.CampaignStatus_Running.String())))

			err := campaignsQuery.Query(cm.app.Db, &runningCampaigns)

			if err != nil {
				cm.app.Logger.Error("error fetching running campaigns from the database", err)
			}

			// ! mutate the active campaign list with the new campaigns, only the newly fetched running campaigns should be in the active campaign list
			for _, campaign := range runningCampaigns {
				if _, ok := cm.activeCampaigns[campaign.Campaign.UniqueId.String()]; !ok {
					cm.activeCampaigns[campaign.Campaign.UniqueId.String()] = struct {
						campaign   model.Campaign
						wapiClient *wapi.Client
					}{
						campaign: campaign.Campaign, wapiClient: wapi.New(&wapi.ClientConfig{
							BusinessAccountId: campaign.WhatsappBusinessAccount.AccountId,
							ApiAccessToken:    campaign.WhatsappBusinessAccount.AccessToken,
							WebhookSecret:     campaign.WhatsappBusinessAccount.WebhookSecret,
							WebhookPath:       "",
						})}
				}
			}

			// ! remove the campaigns that are not running anymore from the active campaign list
			for campaignId, _ := range cm.activeCampaigns {
				// if running campaign list does not have the campaign id, remove it from the active campaign list
				for _, runningCampaign := range runningCampaigns {
					if runningCampaign.Campaign.UniqueId.String() != campaignId {
						delete(cm.activeCampaigns, campaignId)
					}
				}
			}
		}
	}
}

func (cm *CampaignManager) processMessageQueue() {
	for {
		select {
		case message := <-cm.messageQueue:
			err := cm.SendMessage(message)
			if err != nil {
				cm.app.Logger.Error("error sending message to user", err)
				// ! TODO: broadcast this message to websocket via the API server event
			}
		}
	}
}

func (cm *CampaignManager) nextContactsBatch(campaignId string, lastUserId string) {
	var contacts []struct {
		model.Contact
	}

	// ! write a query which creates a sorted list of contacts aggregated from all the campaign list

	// ! this function will be responsible to add the message to the messageQueue of the campaign manager
}

func (cm *CampaignManager) SendMessage(message CampaignMessage) error {
	// client := message.wapiClient
	// messagingClient := client.NewMessagingClient(
	// 	message.phoneNumberToUse,
	// )
	// * create the template message
	// templateMessage, err := wapiComponents.NewTemplateMessage()

	// if err != nil {
	// 	return err
	// }

	// _, err = messagingClient.Message.Send(templateMessage, message.contact.PhoneNumber)

	// if err != nil {
	// 	return err
	// }

	// ! TODO: update the database campaign with the last contact id sent the message to

	return nil
}
