package event_service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"

	cache_service "github.com/wapikit/wapikit/services/redis_service"
)

type EventService struct {
	Db                    *sql.DB
	Logger                *slog.Logger
	Redis                 *cache_service.RedisClient
	RedisEventChannelName string
}

func NewEventService(db *sql.DB, logger *slog.Logger, redis *cache_service.RedisClient, channelName string) *EventService {
	return &EventService{
		Db:                    db,
		Logger:                logger,
		Redis:                 redis,
		RedisEventChannelName: channelName,
	}
}

func (service *EventService) HandleApiServerEvents(ctx context.Context) <-chan ApiServerEventInterface {
	service.Logger.Info("Event service is listening for API server events...")
	streamChannel := make(chan ApiServerEventInterface, 1000)

	fmt.Println("Subscribing to Redis channel", service.RedisEventChannelName)

	redisClient := service.Redis
	pubsub := redisClient.Subscribe(ctx, service.RedisEventChannelName)
	redisEventChannel := pubsub.Channel()

	// Goroutine to listen for Redis events
	go func() {
		defer pubsub.Close()

		for {
			select {
			case apiServerEvent, ok := <-redisEventChannel:
				fmt.Println("API SERVER EVENT RECEIVED")
				if !ok {
					service.Logger.Error("Redis event channel closed, stopping event listener.")
					return
				}

				apiServerEventData := []byte(apiServerEvent.Payload)
				var event BaseApiServerEvent
				err := json.Unmarshal(apiServerEventData, &event)
				if err != nil {
					service.Logger.Error("Unable to unmarshal API server event and determine type", err.Error(), nil)
					continue
				}

				service.Logger.Info("API SERVER EVENT OF TYPE", string(event.EventType), nil)

				switch event.EventType {
				case ApiServerNewMessageEvent:
					var newMessageEvent NewMessageEvent
					err := json.Unmarshal(apiServerEventData, &newMessageEvent)
					if err != nil {
						service.Logger.Error("Unable to unmarshal new message event", err.Error(), nil)
						continue
					}
					streamChannel <- newMessageEvent

				default:
					service.Logger.Info("Unknown event type received")
				}

			case <-ctx.Done(): // Handle context cancellation
				service.Logger.Info("Context cancelled, stopping event stream")
				return
			}
		}
	}()

	return streamChannel
}
