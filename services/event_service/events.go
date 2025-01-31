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
	EventType ApiServerEventType `json:"eventType"`
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
	Conversation ConversationWithAllDetails `json:"conversation"`
	EventType    ApiServerEventType         `json:"eventType"`
	Message      api_types.MessageSchema    `json:"message"`
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
	Conversation ConversationWithAllDetails `json:"conversation"`
	EventType    ApiServerEventType         `json:"eventType"`
	ChatId       string                     `json:"chatId"`
	UserId       string                     `json:"userId"`
}

type ChatUnAssignmentEvent struct {
	BaseApiServerEvent
	EventType ApiServerEventType `json:"eventType"`
	ChatId    string             `json:"chatId"`
	UserId    string             `json:"userId"`
}

type ErrorEvent struct {
	BaseApiServerEvent
	EventType ApiServerEventType `json:"eventType"`
	Error     string             `json:"error"`
}

type ReloadRequiredEvent struct {
	BaseApiServerEvent
	EventType        ApiServerEventType `json:"eventType"`
	IsReloadRequired bool               `json:"isReloadRequired"`
}

type ConversationClosedEvent struct {
	BaseApiServerEvent
	EventType      ApiServerEventType `json:"eventType"`
	ConversationId string             `json:"chatId"`
}
