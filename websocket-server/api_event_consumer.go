package websocket_server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/wapikit/wapi.go/pkg/components"
	"github.com/wapikit/wapikit/internal/core/api_server_events"
	"github.com/wapikit/wapikit/internal/interfaces"
)

// NOTE: we are following a one way data flow for ApiServerEvent, where only the ApiServer itself can update the db for the events changes or any update required like message_log etc.
func HandleApiServerEvents(ctx context.Context, app interfaces.App) {
	redisClient := app.Redis
	pubsub := redisClient.Subscribe(ctx, app.Constants.RedisEventChannelName)
	defer pubsub.Close()
	ch := pubsub.Channel()
	for msg := range ch {
		msgData := []byte(msg.Payload)
		var event api_server_events.BaseApiServerEvent
		err := json.Unmarshal(msgData, &event)
		if err != nil {
			app.Logger.Error("unable to unmarshal api server event and determine type", err.Error(), nil)
			continue
		}

		switch event.EventType {
		case api_server_events.ApiServerChatAssignmentEvent:
			handleChatAssignmentEvent(app)

		case api_server_events.ApiServerNewNotificationEvent:
			handleNewNotificationEvent(app)

		case api_server_events.ApiServerNewMessageEvent:
			handleNewMessageEvent("", app)
		}

		// determine the type of message and call the corresponding handler below
		app.Logger.Info("received message from redis: %v", string(msgData), nil)
	}
}

func handleChatAssignmentEvent(app interfaces.App) {

	// ! get the chat ID, fetch it from the database, then check if the chat is assigned to any connected client, if yes then broadcast the event to them
	// ! wait for the acknowledgment

}

func handleNewNotificationEvent(app interfaces.App) {
	// this will broadcast the notification to the frontend, which is websocketServerEvent

	// get the connection from the connections map
	// send the message to the connection, by building an instance of the WebsocketEventTypeNewNotification
}

// ! TODO:
func handleNewMessageEvent(phoneNumberId string, app interfaces.App) error {
	textMessage, err := components.NewTextMessage(components.TextMessageConfigs{
		Text: "Hii I am websocket message",
	})
	if err != nil {
		return nil
	}

	messagingClient := app.WapiClient.NewMessagingClient(phoneNumberId)
	response, err := messagingClient.Message.Send(textMessage, "919643500545")
	fmt.Println(response)
	// if the response is of error then retry again, but still if the response is error then do send the error at the frontend
	return err
}
