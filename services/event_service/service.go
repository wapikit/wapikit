package event_service

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"

	cache_service "github.com/wapikit/wapikit/services/redis_service"
)

type EventService struct {
	Db                    *sql.DB
	Logger                *slog.Logger
	Redis                 *cache_service.RedisClient
	RedisEventChannelName string
}

func NewEventService(db *sql.DB, logger *slog.Logger, redis *cache_service.RedisClient) *EventService {
	return &EventService{
		Db:     db,
		Logger: logger,
		Redis:  redis,
	}
}

func (service *EventService) HandleApiServerEvents(ctx context.Context) (StreamChannel <-chan string) {
	service.Logger.Info("Event service is listening for api server events...")
	streamChannel := make(chan string)

	redisClient := service.Redis
	pubsub := redisClient.Subscribe(ctx, service.RedisEventChannelName)
	defer pubsub.Close()
	redisEventChannel := pubsub.Channel()

	go func() {
		defer close(streamChannel)

		for apiServerEvent := range redisEventChannel {
			apiServerEventData := []byte(apiServerEvent.Payload)

			var event BaseApiServerEvent
			err := json.Unmarshal(apiServerEventData, &event)
			if err != nil {
				service.Logger.Error("unable to unmarshal api server event and determine type", err.Error(), nil)
				continue
			}

			service.Logger.Info("API SERVER EVENT OF TYPE", string(event.EventType), nil)

			switch event.EventType {

			case ApiServerChatAssignmentEvent:

			case ApiServerNewNotificationEvent:

			case ApiServerNewMessageEvent:
				var event NewMessageEvent
				err := json.Unmarshal(apiServerEventData, &event)
				if err != nil {
					service.Logger.Error("unable to unmarshal new message event", err.Error(), nil)
					continue
				}
				streamChannel <- string(event.ToJson())

			case ApiServerChatUnAssignmentEvent:

			case ApiServerErrorEvent:

			case ApiServerReloadRequiredEvent:

			case ApiServerConversationClosedEvent:

			case ApiServerNewConversationEvent:

			default:
				service.Logger.Info("unknown event type received")
			}
		}
	}()

	return streamChannel
}
