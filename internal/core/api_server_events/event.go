package api_server_events

import (
	"encoding/json"
	"log"

	model "github.com/wapikit/wapikit/.db-generated/model"
	"github.com/wapikit/wapikit/internal/api_types"
)

type ApiServerEventType string

const (
	ApiServerNewNotificationEvent    ApiServerEventType = "NewNotification"
	ApiServerNewMessageEvent         ApiServerEventType = "NewMessage"
	ApiServerChatAssignmentEvent     ApiServerEventType = "ChatAssignment"
	ApiServerChatUnAssignmentEvent   ApiServerEventType = "ChatUnAssignment"
	ApiServerErrorEvent              ApiServerEventType = "Error"
	ApiServerReloadRequiredEvent     ApiServerEventType = "ReloadRequired"
	ApiServerConversationClosedEvent ApiServerEventType = "ConversationClosed"
	ApiServerNewConversationEvent    ApiServerEventType = "NewConversation"
)

type ApiServerEventInterface interface {
	ToJson() []byte
}

type ConversationWithAllDetails struct {
	model.Conversation
	Contact    model.Contact `json:"contact"`
	AssignedTo struct {
		model.OrganizationMember
		User model.User `json:"user"`
	} `json:"assignedTo"`
	WhatsappBusinessAccount struct {
		model.WhatsappBusinessAccount
		Organization model.Organization `json:"organization"`
	} `json:"whatsappBusinessAccount"`
}

type BaseApiServerEvent struct {
	EventType    ApiServerEventType         `json:"eventType"`
	Conversation ConversationWithAllDetails `json:"conversation"`
}

func (event *BaseApiServerEvent) ToJson() []byte {
	bytes, err := json.Marshal(event)
	if err != nil {
		log.Print(err)
	}
	return bytes
}

type NewNotificationEvent struct {
	BaseApiServerEvent                    // make it inline
	EventType          ApiServerEventType `json:"eventType"`
	Notification       string             `json:"notification"`
}

type NewMessageEvent struct {
	BaseApiServerEvent
	EventType ApiServerEventType      `json:"eventType"`
	Message   api_types.MessageSchema `json:"message"`
}

func (event *NewMessageEvent) ToJson() []byte {
	bytes, err := json.Marshal(event)
	if err != nil {
		log.Print(err)
	}
	return bytes
}

type ChatAssignmentEvent struct {
	BaseApiServerEvent
	EventType ApiServerEventType `json:"eventType"`
	ChatId    string             `json:"chatId"`
	UserId    string             `json:"userId"`
}

type ChatUnAssignmentEvent struct {
	BaseApiServerEvent
	EventType ApiServerEventType `json:"eventType"`
	ChatId    string             `json:"chatId"`
	UserId    string             `json:"userId"`
}

// these events are meant to sent to the redis pubsub channel and our websocket server will consume these messages and react to them, also

// ! flow of application:

// ! 1. we have a api server and a websocket server running on two difference thread or in two go routines or may be hosted separately as in case of managed hosted version of the application.
// ! 2. the api server will be handling all the rest api request from the frontend, and if there is something that needs to be immediately sent to the frontend client then the api server will publish a message which in code we wll refer to as ApiServerEvent, the event will then be consumed by the websocket server redis pubsub channel consumer and will be sent to the concerned connection as we have stored an slice of connections in the websocket server.
// ! 3. example of ApiServerEvent can be error event, new notification event, event on a chat assignment to a user, our rest api server is also listening to the whatsapp business webhook so, every time we get a webhook event, and if this is something which requires to be sent on frontend suppose a new message comes in, so we will trigger a event to the redis pubsub channel and the websocket will consume it nd conveys it to the concerned connection.
// ! 4. websocket server is responsible for the tasks mentioned above and whenever there is some message from a app user to a customer, the frontend will send the message over websocket, and then the websocket will call the whatsapp api to send the message to the customer.
