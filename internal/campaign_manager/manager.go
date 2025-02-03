package campaign_manager

import (
	"context"
	"database/sql"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/paulbellamy/ratecounter"
	wapi "github.com/wapikit/wapi.go/pkg/client"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/wapikit/wapikit/.db-generated/model"
	table "github.com/wapikit/wapikit/.db-generated/table"
	"github.com/wapikit/wapikit/api/api_types"
	"github.com/wapikit/wapikit/services/event_service"
	"github.com/wapikit/wapikit/services/notification_service"
	cache_service "github.com/wapikit/wapikit/services/redis_service"

	"github.com/wapikit/wapikit/utils"
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

type businessWorker struct {
	messageQueue chan *CampaignMessage
	rateLimiter  *ratecounter.RateCounter
	wg           sync.WaitGroup
	stopChan     chan struct{}
}

type CampaignManager struct {
	Db     *sql.DB
	Logger slog.Logger

	Redis                 *cache_service.RedisClient
	RedisEventChannelName string

	runningCampaigns      map[string]*runningCampaign
	runningCampaignsMutex sync.RWMutex

	campaignQueue chan *runningCampaign

	businessWorkers      map[string]*businessWorker
	businessWorkersMutex sync.RWMutex

	NotificationService *notification_service.NotificationService
}

func NewCampaignManager(db *sql.DB, logger slog.Logger, redis *cache_service.RedisClient, notificationService *notification_service.NotificationService, redisEventChannelName string) *CampaignManager {
	return &CampaignManager{
		Db:     db,
		Logger: logger,

		runningCampaigns:      make(map[string]*runningCampaign),
		runningCampaignsMutex: sync.RWMutex{},

		// 1000 campaigns can be queued at a time
		campaignQueue: make(chan *runningCampaign, 1000),

		businessWorkers: make(map[string]*businessWorker),

		businessWorkersMutex:  sync.RWMutex{},
		Redis:                 redis,
		RedisEventChannelName: redisEventChannelName,
		NotificationService:   notificationService,
	}
}

type CampaignMessage struct {
	Campaign *runningCampaign `json:"campaign"`
	Contact  model.Contact    `json:"contact"`
}

// New worker function
func (cm *CampaignManager) messageQueueProcessor(businessAccountId string, worker *businessWorker) {
	defer func() {
		// Cleanup when worker stops
		cm.businessWorkersMutex.Lock()
		delete(cm.businessWorkers, businessAccountId)
		cm.businessWorkersMutex.Unlock()
	}()

	for {
		select {
		case <-worker.stopChan:
			return
		case message, ok := <-worker.messageQueue:
			cm.Logger.Debug("sending message", "biz_id", businessAccountId, "campaign_id", message.Campaign.UniqueId.String(), "contact_id", message.Contact.UniqueId.String())

			if !ok {
				return
			}

			if message.Campaign.IsStopped.Load() {
				// * campaign has been stopped, so skip this message
				continue
			}

			// Business-specific rate limiting
			if worker.rateLimiter.Rate() >= int64(messagesPerSecondLimit) {
				// Requeue with backoff
				time.Sleep(10 * time.Millisecond)
				worker.messageQueue <- message
				continue
			}

			cm.sendMessage(message)

			campaignProgressEvent := event_service.NewCampaignProgressEvent(message.Campaign.UniqueId.String(), message.Campaign.Sent.Load(), message.Campaign.ErrorCount.Load(), api_types.Running)
			cm.Redis.PublishMessageToRedisChannel(cm.RedisEventChannelName, campaignProgressEvent.ToJson())
		}
	}
}

func (cm *CampaignManager) newRunningCampaign(dbCampaign model.Campaign, businessAccount model.WhatsappBusinessAccount) *runningCampaign {
	cm.businessWorkersMutex.Lock()
	defer cm.businessWorkersMutex.Unlock()

	cm.Logger.Debug("new campaign started", "campaign_id", dbCampaign.UniqueId.String())

	businessAccountId := businessAccount.AccountId

	if _, exists := cm.businessWorkers[businessAccountId]; !exists {
		worker := &businessWorker{
			messageQueue: make(chan *CampaignMessage, 1000),
			rateLimiter:  ratecounter.NewRateCounter(1 * time.Second),
			stopChan:     make(chan struct{}),
		}

		// Start worker goroutine
		go cm.messageQueueProcessor(businessAccountId, worker)
		cm.businessWorkers[businessAccountId] = worker
	}

	lastContactId := ""

	if dbCampaign.LastContactSent != nil {
		lastContactId = dbCampaign.LastContactSent.String()
	}

	campaign := runningCampaign{
		Campaign: dbCampaign,
		WapiClient: wapi.New(&wapi.ClientConfig{
			BusinessAccountId: businessAccount.AccountId,
			ApiAccessToken:    businessAccount.AccessToken,
			WebhookSecret:     businessAccount.WebhookSecret,
		}),
		PhoneNumberToUse:  dbCampaign.PhoneNumber,
		BusinessAccountId: businessAccount.AccountId,
		LastContactIdSent: lastContactId,
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
	defer cm.Stop()

	// * scan for campaign status changes every 5 seconds
	go cm.queueRunningCampaigns()

	// * scan for scheduled campaign needed to be started every 5 seconds
	go cm.runScheduledCampaigns()

	cm.Logger.Info("campaign manager started.")
	// * process the campaign queue, means listen to the campaign queue, and then for each campaign, call the function to next subscribers
	for campaign := range cm.campaignQueue {
		hasContactsRemainingInQueue := campaign.nextContactsBatch()
		if hasContactsRemainingInQueue {
			cm.Logger.Debug("campaign has contacts remaining in queue", "campaign_id", campaign.UniqueId.String())
			// queue it again
			select {
			case cm.campaignQueue <- campaign:
			default:
			}
		} else {
			cm.Logger.Debug("campaign has no contacts remaining in queue", "campaign_id", campaign.UniqueId.String())
			campaign.wg.Done()
		}
	}
}

func (cm *CampaignManager) updatedCampaignStatus(campaignId uuid.UUID, status model.CampaignStatusEnum) (bool, error) {
	campaignUpdateQuery := table.Campaign.UPDATE(table.Campaign.Status).
		SET(status).
		WHERE(table.Campaign.UniqueId.EQ(UUID(campaignId)))

	_, err := campaignUpdateQuery.Exec(cm.Db)

	if err != nil {
		cm.Logger.Error("error updating campaign status", err.Error())
		return false, err
	}

	return true, nil
}

func (cm *CampaignManager) queueRunningCampaigns() {
	// * scan for campaign status changes every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			currentRunningCampaignIds := cm.getRunningCampaignsUniqueIds()

			cm.Logger.Debug("running campaigns", "campaigns", currentRunningCampaignIds)

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

			// * if there are running campaigns already in progress, ignore them to be fetched again from DB
			if len(runningCampaignExpression) > 0 {
				whereCondition = whereCondition.AND(
					table.Campaign.UniqueId.NOT_IN(runningCampaignExpression...),
				)
			}

			var runningCampaigns []struct {
				model.Campaign
				model.WhatsappBusinessAccount
			}

			campaignsQuery := SELECT(table.Campaign.AllColumns, table.WhatsappBusinessAccount.AllColumns).
				FROM(table.Campaign.LEFT_JOIN(table.WhatsappBusinessAccount, table.WhatsappBusinessAccount.OrganizationId.EQ(table.Campaign.OrganizationId))).
				WHERE(whereCondition)

			context := context.Background()
			err := campaignsQuery.QueryContext(context, cm.Db, &runningCampaigns)

			if err != nil {
				cm.Logger.Error("error fetching running campaigns from the database", err)
			}

			if len(runningCampaigns) == 0 {
				// no running campaign found
				cm.Logger.Debug("no running campaigns found")
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

func (cm *CampaignManager) runScheduledCampaigns() {
	// * scan for scheduled campaign needed to be started every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			var scheduledCampaigns []model.Campaign

			campaignsQuery := SELECT(table.Campaign.AllColumns).
				FROM(table.Campaign).
				WHERE(table.Campaign.Status.EQ(utils.EnumExpression(model.CampaignStatusEnum_Scheduled.String())))

			context := context.Background()
			err := campaignsQuery.QueryContext(context, cm.Db, &scheduledCampaigns)

			if err != nil {
				cm.Logger.Error("error fetching scheduled campaigns from the database", err)
			}

			if len(scheduledCampaigns) == 0 {
				// no scheduled campaign found
				continue
			}

			for _, campaign := range scheduledCampaigns {
				// * check if the scheduled time has passed, if yes then update the campaign status to running
				if campaign.ScheduledAt.Before(time.Now()) {
					_, err := cm.updatedCampaignStatus(campaign.UniqueId, model.CampaignStatusEnum_Running)
					if err != nil {
						cm.Logger.Error("error updating campaign status to running", err.Error())
					}
				}
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

func (cm *CampaignManager) UpdateLastContactId(campaignId, lastContactId uuid.UUID) error {
	campaignUpdateQuery := table.Campaign.UPDATE(table.Campaign.LastContactSent).
		SET(lastContactId).
		WHERE(table.Campaign.UniqueId.EQ(UUID(campaignId)))

	_, err := campaignUpdateQuery.Exec(cm.Db)

	if err != nil {
		cm.Logger.Error("error updating campaign last contact id", err.Error())
		return err
	}

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

func (cm *CampaignManager) Stop() {
	cm.businessWorkersMutex.Lock()
	defer cm.businessWorkersMutex.Unlock()

	for businessAccountId, worker := range cm.businessWorkers {
		close(worker.stopChan)
		close(worker.messageQueue)
		delete(cm.businessWorkers, businessAccountId)
	}
}
