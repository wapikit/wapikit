package campaign_manager

import (
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
	wapi "github.com/wapikit/wapi.go/pkg/client"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/wapikit/wapikit/.db-generated/model"
	table "github.com/wapikit/wapikit/.db-generated/table"
	"github.com/wapikit/wapikit/api/api_types"
	"github.com/wapikit/wapikit/services/event_service"
	"github.com/wapikit/wapikit/services/notification_service"
)

type runningCampaign struct {
	model.Campaign
	WapiClient        *wapi.Client `json:"wapiclient"`
	BusinessAccountId string       `json:"businessAccountId"`
	PhoneNumberToUse  string       `json:"phoneNumberToUse"`

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
	rc.Manager.Logger.Info("fetching next contact batch", nil)
	var contacts []model.Contact

	contactsCte := CTE("contacts")

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

		// Get business worker
		rc.Manager.businessWorkersMutex.RLock()
		worker, exists := rc.Manager.businessWorkers[rc.BusinessAccountId]
		rc.Manager.businessWorkersMutex.RUnlock()

		if !exists {
			rc.Manager.Logger.Error("Business worker not found", nil)
			rc.Manager.NotificationService.SendSlackNotification(notification_service.SlackNotificationParams{
				Title:   "🚨🚨 Business worker not found in next contact batch 🚨🚨",
				Message: "Business worker not found for business account ID: " + rc.BusinessAccountId,
			})

			continue
		}

		select {
		case worker.messageQueue <- message:
			rc.Manager.Logger.Info("message added to the message queue from next contacts batch", nil)
			rc.wg.Add(1)
		default:
			// * if the message queue is full, then return true, so that the campaign can be queued again
			return true
		}
	}

	rc.LastContactIdSent = contacts[len(contacts)-1].UniqueId.String()
	return len(contacts) > 0
}

func (rc *runningCampaign) stop() {
	rc.Manager.Logger.Info("stopping the campaign", nil)
	if rc.IsStopped.Load() {
		return
	}
	rc.IsStopped.Store(true)
}

// this function will only run when the campaign is exhausted its subscriber list
func (rc *runningCampaign) cleanUp() {

	rc.Manager.Logger.Info("cleaning up the campaign", nil)

	defer func() {
		rc.Manager.runningCampaignsMutex.Lock()
		delete(rc.Manager.runningCampaigns, rc.UniqueId.String())
		rc.Manager.runningCampaignsMutex.Unlock()
	}()

	// check the fresh status of the campaign, if it is still running, then update the status to finished
	var campaign model.Campaign

	campaignQuery := SELECT(table.Campaign.AllColumns).
		FROM(table.Campaign).
		WHERE(table.Campaign.UniqueId.EQ(UUID(rc.UniqueId)))

	err := campaignQuery.Query(rc.Manager.Db, &campaign)

	if err != nil {
		rc.Manager.Logger.Error("error fetching campaign from the database", err.Error(), nil)
		// campaign not found in the db for some reason, it will be removed from the running campaigns list
		return
	}

	if campaign.Status == model.CampaignStatusEnum_Running {
		_, err = rc.Manager.updatedCampaignStatus(rc.UniqueId, model.CampaignStatusEnum_Finished)
		if err != nil {
			rc.Manager.Logger.Error("error updating campaign status", err.Error(), nil)
		}
	}

	campaignProgressEvent := event_service.NewCampaignProgressEvent(rc.UniqueId.String(), rc.Sent.Load(), rc.ErrorCount.Load(), api_types.Finished)
	err = rc.Manager.Redis.PublishMessageToRedisChannel(rc.Manager.RedisEventChannelName, campaignProgressEvent.ToJson())
}
