package websocket_server

import (
	"encoding/json"
	"fmt"
)

// * these are the event send to and from connected clients

type WebsocketEventType string

const (
	WebsocketEventTypeMessageAcknowledgement WebsocketEventType = "MessageAcknowledgementEvent"
	WebsocketEventTypeMessage                WebsocketEventType = "MessageEvent"
	WebsocketEventTypeNotificationRead       WebsocketEventType = "NotificationReadEvent"
	WebsocketEventTypeMessageRead            WebsocketEventType = "MessageReadEvent"
	WebsocketEventTypeNewNotification        WebsocketEventType = "NewNotificationEvent"
	WebsocketEventTypeSystemReload           WebsocketEventType = "SystemReloadRequired"
	WebsocketEventTypeConversationAssignment WebsocketEventType = "ConversationAssignmentEvent"
	WebsocketEventTypeConversationClosed     WebsocketEventType = "ConversationClosedEvent"
	WebsocketEventTypeNewConversation        WebsocketEventType = "NewConversationEvent"
	WebsocketEventTypePing                   WebsocketEventType = "PingEvent"
)

type WebsocketEvent struct {
	EventName WebsocketEventType `json:"eventName"`
	Data      json.RawMessage    `json:"data"`
	EventId   string             `json:"eventId"`
}

func (event WebsocketEvent) toJson() []byte {
	jsonMessage, err := json.Marshal(event)
	if err != nil {
		fmt.Errorf("Error occurred while converting data to json")
	}
	return jsonMessage
}

type WebsocketEventDataInterface interface {
	GetEventName() string
}

type BaseWebsocketEventData struct {
	EventName WebsocketEventType `json:"eventName"`
}

func (event BaseWebsocketEventData) GetEventName() string {
	return string(event.EventName)
}

type MessageAcknowledgementEventData struct {
	Message string `json:"message"`
}

func NewAcknowledgementEvent(eventId string, message string) *WebsocketEvent {
	data := MessageAcknowledgementEventData{
		Message: message,
	}
	marshalData, err := json.Marshal(data)
	if err != nil {
		fmt.Errorf("Error occurred while converting data to json")
	}
	return &WebsocketEvent{
		EventName: WebsocketEventTypeMessageAcknowledgement,
		Data:      marshalData,
		EventId:   eventId,
	}
}

type PingEventData struct {
	Data string `json:"data"`
}

type MessageEventData struct {
	BaseWebsocketEventData `json:"-,inline"`
	Data                   struct {
		MessageID      string `json:"messageId"`
		ConversationID string `json:"conversationId"`
		Message        string `json:"message"`
		SenderID       string `json:"senderId"`
		SenderName     string `json:"senderName"`
		SenderAvatar   string `json:"senderAvatar"`
		SentAt         string `json:"sentAt"`
		IsRead         bool   `json:"isRead"`
	} `json:"data"`
}

type NotificationReadEventData struct {
	BaseWebsocketEventData `json:"-,inline"`
	Data                   struct {
		NotificationID string `json:"notificationId"`
	} `json:"data"`
}

type MessageReadEventData struct {
	BaseWebsocketEventData `json:"-,inline"`
	Data                   struct {
		MessageID string `json:"messageId"`
	} `json:"data"`
}

type NewNotificationEventData struct {
	BaseWebsocketEventData `json:"-,inline"`
	Data                   struct {
		NotificationID string `json:"notificationId"`
	} `json:"data"`
}

type SystemReloadEventData struct {
	BaseWebsocketEventData `json:"-,inline"`
	Data                   struct {
		MessageText      string `json:"messageText"`
		MessageTitle     string `json:"messageTitle"`
		IsReloadRequired bool   `json:"isReloadRequired"`
	} `json:"data"`
}

type ConversationAssignmentEventData struct {
	BaseWebsocketEventData `json:"-,inline"`
	Data                   struct {
		AssignedToMemberID string `json:"assignedToMemberId"`
		ConversationID     string `json:"conversationId"`
		AssignedAt         string `json:"assignedAt"`
	} `json:"data"`
}

type ConversationClosedEventData struct {
	BaseWebsocketEventData `json:"-,inline"`
	Data                   struct {
		ConversationID string `json:"conversationId"`
	} `json:"data"`
}

type NewConversationEventData struct {
	BaseWebsocketEventData `json:"-,inline"`
	Data                   struct {
		ConversationID string `json:"conversationId"`
	} `json:"data"`
}
