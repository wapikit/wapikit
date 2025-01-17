package websocket_server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/wapikit/wapikit/internal/api_server_events"
	"github.com/wapikit/wapikit/internal/interfaces"
	"github.com/wapikit/wapikit/internal/utils"
)

// * these are the event handlers for the api server events, which are published by the api server and consumed by the websocket server

// NOTE: we are following a one way data flow for ApiServerEvent, where only the ApiServer itself can update the db for the events changes or any update required like message_log etc.
func (server *WebSocketServer) HandleApiServerEvents(ctx context.Context, app interfaces.App) {
	app.Logger.Info("websocket server is listening for api server events...")

	redisClient := app.Redis
	pubsub := redisClient.Subscribe(ctx, app.Constants.RedisEventChannelName)
	defer pubsub.Close()
	ch := pubsub.Channel()

	for apiServerEvent := range ch {
		apiServerEventData := []byte(apiServerEvent.Payload)

		var event api_server_events.BaseApiServerEvent
		err := json.Unmarshal(apiServerEventData, &event)
		if err != nil {
			app.Logger.Error("unable to unmarshal api server event and determine type", err.Error(), nil)
			continue
		}

		app.Logger.Info("API SERVER EVENT OF TYPE", string(event.EventType), nil)

		switch event.EventType {

		case api_server_events.ApiServerChatAssignmentEvent:
			handleChatAssignmentEvent(app)

		case api_server_events.ApiServerNewNotificationEvent:
			handleNewNotificationEvent(app)

		case api_server_events.ApiServerNewMessageEvent:
			var event api_server_events.NewMessageEvent
			err := json.Unmarshal(apiServerEventData, &event)
			if err != nil {
				app.Logger.Error("unable to unmarshal new message event", err.Error(), nil)
				continue
			}
			handleNewMessageEvent(app, *server, event)

		case api_server_events.ApiServerChatUnAssignmentEvent:

		case api_server_events.ApiServerErrorEvent:

		case api_server_events.ApiServerReloadRequiredEvent:

		case api_server_events.ApiServerConversationClosedEvent:

		case api_server_events.ApiServerNewConversationEvent:

		default:
			app.Logger.Info("unknown event type received")
		}

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

func handleNewMessageEvent(app interfaces.App, ws WebSocketServer, event api_server_events.NewMessageEvent) error {
	// * this event means we have received a new message from the whatsapp webhook, so we have to broadcast it to the frontend client if connected
	fmt.Println("websocket server have received a new message to broadcast to the frontend", event.Message)
	newMessageReceivedWebsocketEvent := NewMessageReceivedWebsocketEvent(utils.GenerateWebsocketEventId(), event.Message)
	fmt.Println("sending message to client", string(newMessageReceivedWebsocketEvent.toJson()))
	errors := ws.broadcastToAll(newMessageReceivedWebsocketEvent.toJson())

	if len(errors) > 0 {
		app.Logger.Error("error sending message to client", errors, nil)
		return fmt.Errorf("error sending message to clients")
		// ! TODO retry the message sending
	}

	return nil
}
