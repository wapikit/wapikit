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

			cm.Logger.Info("sending message", "biz_id", businessAccountId, "campaign_id", message.Campaign.UniqueId.String(), "contact_id", message.Contact.UniqueId.String())

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

			err := cm.sendMessage(message)

			campaignProgressEvent := event_service.NewCampaignProgressEvent(message.Campaign.UniqueId.String(), message.Campaign.Sent.Load(), message.Campaign.ErrorCount.Load(), api_types.Running)
			err = cm.Redis.PublishMessageToRedisChannel(cm.RedisEventChannelName, campaignProgressEvent.ToJson())

			if err != nil {
				cm.Logger.Error("error sending message", "biz_id", businessAccountId, "error", err)
			}

			worker.rateLimiter.Incr(1)
		}
	}
}

func (cm *CampaignManager) newRunningCampaign(dbCampaign model.Campaign, businessAccount model.WhatsappBusinessAccount) *runningCampaign {
	cm.businessWorkersMutex.Lock()
	defer cm.businessWorkersMutex.Unlock()

	cm.Logger.Info("new campaign started", "campaign_id", dbCampaign.UniqueId.String())

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
	// * scan for campaign status changes every 5 seconds
	go cm.queueRunningCampaigns()

	// * scan for scheduled campaign needed to be started every 5 seconds
	go cm.runScheduledCampaigns()

	cm.Logger.Info("campaign manager started.")
	// * process the campaign queue, means listen to the campaign queue, and then for each campaign, call the function to next subscribers
	for campaign := range cm.campaignQueue {
		hasContactsRemainingInQueue := campaign.nextContactsBatch()
		if hasContactsRemainingInQueue {
			cm.Logger.Info("campaign has contacts remaining in queue", "campaign_id", campaign.UniqueId.String())
			// queue it again
			select {
			case cm.campaignQueue <- campaign:
			default:
			}
		} else {
			cm.Logger.Info("campaign has no contacts remaining in queue", "campaign_id", campaign.UniqueId.String())
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

			cm.Logger.Info("running campaigns", "campaigns", currentRunningCampaignIds)

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
				cm.Logger.Info("no running campaigns found")
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

func (cm *CampaignManager) sendMessage(message *CampaignMessage) error {
	// * irrespective of the error, we need to decrement the wait group, because the message has been tried to be sent
	// ! TODO: add a retry mechanism in the future and also store the campaign logs in the db
	defer message.Campaign.wg.Done()
	defer func() {
		err := cm.UpdateLastContactId(message.Campaign.UniqueId, message.Contact.UniqueId)
		if err != nil {
			cm.Logger.Error("error updating last contact id", err.Error())
		}
	}()

	// Get business worker
	cm.businessWorkersMutex.RLock()
	worker, exists := cm.businessWorkers[message.Campaign.BusinessAccountId]
	cm.businessWorkersMutex.RUnlock()

	if !exists {
		// Handle worker not found (shouldn't happen)
		cm.Logger.Error("Business worker not found", nil)

		cm.NotificationService.SendSlackNotification(notification_service.SlackNotificationParams{
			Title:   "🚨🚨 Business worker not found in send message 🚨🚨",
			Message: "Business worker not found for business account ID: " + message.Campaign.BusinessAccountId,
		})

		return fmt.Errorf("business worker not found for business account ID: %s", message.Campaign.BusinessAccountId)
	}

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
							if len(parameterStoredInDb.Buttons) > 0 {
								templateMessage.AddButton(wapiComponents.TemplateMessageComponentButtonType{
									Type:    wapiComponents.TemplateMessageComponentTypeButton,
									SubType: wapiComponents.TemplateMessageButtonComponentTypeUrl,
									Index:   index,
								})
							} else {
								templateMessage.AddButton(wapiComponents.TemplateMessageComponentButtonType{
									Type:       wapiComponents.TemplateMessageComponentTypeButton,
									SubType:    wapiComponents.TemplateMessageButtonComponentTypeUrl,
									Index:      index,
									Parameters: []wapiComponents.TemplateMessageParameter{}},
								)
							}

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

	messageStatus := model.MessageStatusEnum_Sent

	if err != nil {
		fmt.Errorf("error sending message to user: %v", err.Error())
		message.Campaign.ErrorCount.Add(1)
		messageStatus = model.MessageStatusEnum_Failed
		return err
	}

	jsonMessage, err := templateMessage.ToJson(wapiComponents.ApiCompatibleJsonConverterConfigs{
		SendToPhoneNumber: message.Contact.PhoneNumber,
	})

	stringifiedJsonMessage := string(jsonMessage)

	if err != nil {
		return err
	} else {
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
			Status:          messageStatus,
		}

		messageSentRecordQuery := table.Message.
			INSERT(table.Message.MutableColumns).
			MODEL(messageSent).
			RETURNING(table.Message.AllColumns)

		err := messageSentRecordQuery.Query(cm.Db, &messageSent)

		if err != nil {
			cm.Logger.Error("error saving message record to the database", err.Error())
		}
	}

	worker.rateLimiter.Incr(1)
	message.Campaign.Sent.Add(1)
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

// Add cleanup to CampaignManager
func (cm *CampaignManager) Stop() {
	cm.businessWorkersMutex.Lock()
	defer cm.businessWorkersMutex.Unlock()

	for businessAccountId, worker := range cm.businessWorkers {
		close(worker.stopChan)
		close(worker.messageQueue)
		delete(cm.businessWorkers, businessAccountId)
	}
}
