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
	ApiServerCampaignProgressEvent   ApiServerEventType = "CampaignProgress"
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
	EventType      ApiServerEventType `json:"event"`
	Data           interface{}        `json:"data"`
	UserId         string             `json:"userId"`
	OrganizationId string             `json:"organizationId"`
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

// Specific event types with properly structured Data

type NewNotificationEvent struct {
	BaseApiServerEvent
}

func NewNewNotificationEvent(notification string) *NewNotificationEvent {
	return &NewNotificationEvent{
		BaseApiServerEvent: BaseApiServerEvent{
			EventType: ApiServerNewNotificationEvent,
			Data: struct {
				Notification string `json:"notification"`
			}{
				Notification: notification,
			},
		},
	}
}

type NewMessageEvent struct {
	BaseApiServerEvent
}

func NewNewMessageEvent(conv ConversationWithAllDetails, msg api_types.MessageSchema, userId, orgId string) *NewMessageEvent {
	return &NewMessageEvent{
		BaseApiServerEvent: BaseApiServerEvent{
			EventType: ApiServerNewMessageEvent,
			Data: struct {
				Conversation ConversationWithAllDetails `json:"conversation"`
				Message      api_types.MessageSchema    `json:"message"`
			}{
				Conversation: conv,
				Message:      msg,
			},
			UserId:         userId,
			OrganizationId: orgId,
		},
	}
}

type ChatAssignmentEvent struct {
	BaseApiServerEvent
}

func NewChatAssignmentEvent(conv ConversationWithAllDetails, chatId, userId string, orgId string) *ChatAssignmentEvent {
	return &ChatAssignmentEvent{
		BaseApiServerEvent: BaseApiServerEvent{
			EventType: ApiServerChatAssignmentEvent,
			Data: struct {
				Conversation ConversationWithAllDetails `json:"conversation"`
				ChatId       string                     `json:"chatId"`
				UserId       string                     `json:"userId"`
			}{
				Conversation: conv,
				ChatId:       chatId,
				UserId:       userId,
			},
			UserId:         userId,
			OrganizationId: orgId,
		},
	}
}

type ChatUnAssignmentEvent struct {
	BaseApiServerEvent
}

func NewChatUnAssignmentEvent(chatId, userId string) *ChatUnAssignmentEvent {
	return &ChatUnAssignmentEvent{
		BaseApiServerEvent: BaseApiServerEvent{
			EventType: ApiServerChatUnAssignmentEvent,
			Data: struct {
				ChatId string `json:"chatId"`
				UserId string `json:"userId"`
			}{
				ChatId: chatId,
				UserId: userId,
			},
		},
	}
}

type ErrorEvent struct {
	BaseApiServerEvent
}

func NewErrorEvent(errorMsg string) *ErrorEvent {
	return &ErrorEvent{
		BaseApiServerEvent: BaseApiServerEvent{
			EventType: ApiServerErrorEvent,
			Data: struct {
				Error string `json:"error"`
			}{
				Error: errorMsg,
			},
		},
	}
}

type ReloadRequiredEvent struct {
	BaseApiServerEvent
}

func NewReloadRequiredEvent(isReloadRequired bool) *ReloadRequiredEvent {
	return &ReloadRequiredEvent{
		BaseApiServerEvent: BaseApiServerEvent{
			EventType: ApiServerReloadRequiredEvent,
			Data: struct {
				IsReloadRequired bool `json:"isReloadRequired"`
			}{
				IsReloadRequired: isReloadRequired,
			},
		},
	}
}

type ConversationClosedEvent struct {
	BaseApiServerEvent
}

func NewConversationClosedEvent(conversationId string) *ConversationClosedEvent {
	return &ConversationClosedEvent{
		BaseApiServerEvent: BaseApiServerEvent{
			EventType: ApiServerConversationClosedEvent,
			Data: struct {
				ConversationId string `json:"chatId"`
			}{
				ConversationId: conversationId,
			},
		},
	}
}

type NewConversationEvent struct {
	BaseApiServerEvent
}

func NewNewConversationEvent(conv ConversationWithAllDetails) *NewConversationEvent {
	return &NewConversationEvent{
		BaseApiServerEvent: BaseApiServerEvent{
			EventType: ApiServerNewConversationEvent,
			Data: struct {
				Conversation ConversationWithAllDetails `json:"conversation"`
			}{
				Conversation: conv,
			},
		},
	}
}

type CampaignProgressEventData struct {
	CampaignId      string                       `json:"campaignId"`
	MessagesSent    int64                        `json:"messagesSent"`
	MessagesErrored int64                        `json:"messagesErrored"`
	Status          api_types.CampaignStatusEnum `json:"status"`
}

type CampaignProgressEvent struct {
	BaseApiServerEvent
}

func NewCampaignProgressEvent(campaignId string, messagesSent, messagesErrored int64, status api_types.CampaignStatusEnum) *CampaignProgressEvent {
	return &CampaignProgressEvent{
		BaseApiServerEvent: BaseApiServerEvent{
			EventType: ApiServerCampaignProgressEvent,
			Data: CampaignProgressEventData{
				CampaignId:      campaignId,
				MessagesSent:    messagesSent,
				MessagesErrored: messagesErrored,
				Status:          status,
			},
		},
	}
}
