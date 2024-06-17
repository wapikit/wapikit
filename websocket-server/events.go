package websocket_server

import (
	"encoding/json"
	"fmt"
)

type WebsocketEventType string

const (
	WebsocketEventTypeMessageAcknowledgement WebsocketEventType = "MessageAcknowledgement"
	WebsocketEventTypeMessage                WebsocketEventType = "MessageEvent"
	WebsocketEventTypeNotificationRead       WebsocketEventType = "notification_read"
	WebsocketEventTypeMessageRead            WebsocketEventType = "message_read"
	WebsocketEventTypeNewNotification        WebsocketEventType = "new_notification"
	WebsocketEventTypeSystemReload           WebsocketEventType = "system_reload"
	WebsocketEventTypeConversationAssignment WebsocketEventType = "conversation_assignment"
	WebsocketEventTypeConversationClosed     WebsocketEventType = "conversation_closed"
	WebsocketEventTypeNewConversation        WebsocketEventType = "new_conversation"
	WebsocketEventTypePing                   WebsocketEventType = "ping"
)

type WebsocketEvent struct {
	EventName WebsocketEventType          `json:"eventName"`
	Data      WebsocketEventDataInterface `json:"data"`
	MessageId string                      `json:"messageId"`
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
	BaseWebsocketEventData `json:"-,inline"`
	Data                   struct {
		MessageID string `json:"messageId"`
		Message   string `json:"message"`
	}
}

func NewAcknowledgementEvent(messageId string, message string) *WebsocketEvent {
	return &WebsocketEvent{
		EventName: WebsocketEventTypeMessageAcknowledgement,
		Data: MessageAcknowledgementEventData{
			Data: struct {
				MessageID string "json:\"messageId\""
				Message   string "json:\"message\""
			}{
				MessageID: messageId,
				Message:   message,
			}},
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
