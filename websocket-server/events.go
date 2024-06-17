package websocket_server

type WebsocketEvent struct {
	Type string                      `json:"type"`
	Data WebsocketEventDataInterface `json:"data"`
}

type WebsocketEventDataInterface interface {
	toJson() string
}

type BaseWebsocketEventData struct {
	EventName string `json:"eventName"`
}

func (b BaseWebsocketEventData) toJson() string {
	return ""
}

type MessageEventData struct {
	EventName string `json:"eventName"`
	Data      struct {
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
	EventName string `json:"eventName"`
	Data      struct {
		NotificationID string `json:"notificationId"`
	} `json:"data"`
}

type MessageReadEventData struct {
	EventName string `json:"eventName"`
	Data      struct {
		MessageID string `json:"messageId"`
	} `json:"data"`
}

type NewNotificationEventData struct {
	EventName string `json:"eventName"`
	Data      struct {
		NotificationID string `json:"notificationId"`
	} `json:"data"`
}

type SystemReloadEventData struct {
	EventName string `json:"eventName"`
	Data      struct {
		MessageText      string `json:"messageText"`
		MessageTitle     string `json:"messageTitle"`
		IsReloadRequired bool   `json:"isReloadRequired"`
	} `json:"data"`
}

type ConversationAssignmentEventData struct {
	EventName string `json:"eventName"`
	Data      struct {
		AssignedToMemberID string `json:"assignedToMemberId"`
		ConversationID     string `json:"conversationId"`
		AssignedAt         string `json:"assignedAt"`
	} `json:"data"`
}

type ConversationClosedEventData struct {
	EventName string `json:"eventName"`
	Data      struct {
		ConversationID string `json:"conversationId"`
	} `json:"data"`
}

type NewConversationEventData struct {
	EventName string `json:"eventName"`
	Data      struct {
		ConversationID string `json:"conversationId"`
	} `json:"data"`
}
