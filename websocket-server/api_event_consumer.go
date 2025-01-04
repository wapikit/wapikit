package websocket_server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/wapikit/wapikit/internal/core/api_server_events"
	"github.com/wapikit/wapikit/internal/interfaces"
)

// * these are the event handlers for the api server events, which are published by the api server and consumed by the websocket server

// NOTE: we are following a one way data flow for ApiServerEvent, where only the ApiServer itself can update the db for the events changes or any update required like message_log etc.
func (server *WebSocketServer) HandleApiServerEvents(ctx context.Context, app interfaces.App) {
	app.Logger.Info("websocket server is listening for api server events...")
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

		fmt.Println("API SERVER EVENT OF TYPE", event.EventType)
		switch event.EventType {
		case api_server_events.ApiServerChatAssignmentEvent:
			handleChatAssignmentEvent(app)

		case api_server_events.ApiServerNewNotificationEvent:
			handleNewNotificationEvent(app)

		case api_server_events.ApiServerNewMessageEvent:
			var event api_server_events.NewMessageEvent
			err := json.Unmarshal(msgData, &event)
			if err != nil {
				fmt.Println("error unmarshalling new message event", err.Error())
				app.Logger.Error("unable to unmarshal new message event", err.Error(), nil)
				continue
			}
			handleNewMessageEvent(app, *server, event)
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
func handleNewMessageEvent(app interfaces.App, ws WebSocketServer, event api_server_events.NewMessageEvent) error {
	// * this event means we have received a new message from the whatsapp webhook, so we have to broadcast it to the frontend client if connected

	fmt.Println("websocket server have received a new message to broadcast to the frontend", event.Message)

	// get the connection from the connections map

	// get the first connection from the connections map

	var conn *WebsocketConnectionData

	for _, connection := range ws.connections {
		conn = connection
		break
	}

	if conn == nil {
		app.Logger.Info("no connection found to broadcast the message")
		return nil
	}

	err := ws.sendMessageToClient(conn.Connection, []byte(event.Message))

	if err != nil {
		fmt.Println("error sending message to client", err.Error())
		// ! retry the message sending
	}

	// textMessage, err := components.NewTextMessage(components.TextMessageConfigs{
	// 	Text: "Hii I am websocket message",
	// })
	// if err != nil {
	// 	return nil
	// }

	// if the response is of error then retry again, but still if the response is error then do send the error at the frontend
	return nil
}
