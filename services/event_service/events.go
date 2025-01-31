package event_service

import (
	"encoding/json"
	"log"

	model "github.com/wapikit/wapikit/.db-generated/model"
	"github.com/wapikit/wapikit/api/api_types"
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
	GetEventType() ApiServerEventType
	GetData() interface{}
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
	EventType ApiServerEventType `json:"event"`
	Data      interface{}        `json:"data"`
}

func (event BaseApiServerEvent) ToJson() []byte {
	bytes, err := json.Marshal(event)
	if err != nil {
		log.Print(err)
	}
	return bytes
}

func (event BaseApiServerEvent) GetEventType() ApiServerEventType {
	return event.EventType
}

func (event BaseApiServerEvent) GetData() interface{} {
	return event.Data
}

type NewNotificationEvent struct {
	BaseApiServerEvent                    // make it inline
	EventType          ApiServerEventType `json:"eventType"`
	Data               struct {
		Notification string `json:"notification"`
	} `json:"data"`
}

type NewMessageEvent struct {
	BaseApiServerEvent
	EventType ApiServerEventType `json:"event"`
	Data      struct {
		Conversation ConversationWithAllDetails `json:"conversation"`
		Message      api_types.MessageSchema    `json:"message"`
	} `json:"data"`
}

func (event NewMessageEvent) ToJson() []byte {
	bytes, err := json.Marshal(event)
	if err != nil {
		log.Print(err)
	}
	return bytes
}

func (event NewMessageEvent) GetEventType() ApiServerEventType {
	return event.EventType
}

func (event NewMessageEvent) GetData() interface{} {
	return event.Data
}

type ChatAssignmentEvent struct {
	BaseApiServerEvent
	Conversation ConversationWithAllDetails `json:"conversation"`
	EventType    ApiServerEventType         `json:"event"`
	ChatId       string                     `json:"chatId"`
	UserId       string                     `json:"userId"`
}

type ChatUnAssignmentEvent struct {
	BaseApiServerEvent
	EventType ApiServerEventType `json:"event"`
	ChatId    string             `json:"chatId"`
	UserId    string             `json:"userId"`
}

type ErrorEvent struct {
	BaseApiServerEvent
	EventType ApiServerEventType `json:"event"`
	Error     string             `json:"error"`
}

type ReloadRequiredEvent struct {
	BaseApiServerEvent
	EventType        ApiServerEventType `json:"event"`
	IsReloadRequired bool               `json:"isReloadRequired"`
}

type ConversationClosedEvent struct {
	BaseApiServerEvent
	EventType      ApiServerEventType `json:"event"`
	ConversationId string             `json:"chatId"`
}
