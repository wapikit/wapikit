package websocket_server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sarthakjdev/wapikit/internal"
	"github.com/sarthakjdev/wapikit/internal/api_server_events"
	"github.com/sarthakjdev/wapikit/internal/interfaces"

	"github.com/sarthakjdev/wapi.go/pkg/components"
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
			app.Logger.Error("Unable to unmarshal api server event and determine type", err.Error())
			continue
		}

		switch event.EventType {
		case api_server_events.ApiServerChatAssignmentEvent:
			handleChatAssignmentEvent()

		case api_server_events.ApiServerNewNotificationEvent:
			handleNewNotificationEvent()

		case api_server_events.ApiServerNewMessageEvent:
			handleNewMessageEvent(app)
		}

		// determine the type of message and call the corresponding handler below
		app.Logger.Info("received message from redis: %v", string(msgData))
	}
}

func handleChatAssignmentEvent() {}

func handleNewNotificationEvent() {
	// this will broadcast the notification to the frontend, which is websocketServerEvent

	// get the connection from the connections map
	// send the message to the connection, by building an instance of the WebsocketEventTypeNewNotification
}

func handleNewMessageEvent(app interfaces.App) error {
	wapiClient := internal.GetWapiClient(&app)
	textMessage, err := components.NewTextMessage(components.TextMessageConfigs{
		Text: "Hii I am websocket message",
	})
	if err != nil {
		return nil
	}
	response, err := wapiClient.Message.Send(textMessage, "919643500545")
	fmt.Println(response)
	// if the response is of error then retry again, but still if the response is error then do send the error at the frontend
	return err
}
