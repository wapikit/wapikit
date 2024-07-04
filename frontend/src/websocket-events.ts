import { z } from 'zod'

export enum WebsocketEventEnum {
	MessageAcknowledgementEvent = 'MessageAcknowledgementEvent',
	MessageEvent = 'MessageEvent',
	NotificationReadEvent = 'NotificationReadEvent',
	MessageReadEvent = 'MessageReadEvent',
	NewNotificationEvent = 'NewNotificationEvent',
	SystemReloadEvent = 'SystemReloadEvent',
	ConversationAssignmentEvent = 'ConversationAssignmentEvent',
	ConversationClosedEvent = 'ConversationClosedEvent',
	NewConversationEvent = 'NewConversationEvent',
	PingEvent = 'PingEvent'
}

export const WebsocketEventDataMap = {
	[WebsocketEventEnum.MessageEvent]: z.object({
		eventName: z.literal(WebsocketEventEnum.MessageEvent),
		messageId: z.string(),
		data: z.object({
			messageId: z.string(),
			conversationId: z.string(),
			message: z.string(),
			senderId: z.string(),
			senderName: z.string(),
			senderAvatar: z.string(),
			sentAt: z.string(),
			isRead: z.boolean()
		})
	}),
	[WebsocketEventEnum.NotificationReadEvent]: z.object({
		eventName: z.literal(WebsocketEventEnum.NotificationReadEvent),
		messageId: z.string(),
		data: z.object({
			notificationId: z.string()
		})
	}),
	[WebsocketEventEnum.MessageReadEvent]: z.object({
		messageId: z.string(),
		eventName: z.literal(WebsocketEventEnum.MessageReadEvent),
		data: z.object({
			messageId: z.string()
		})
	}),
	[WebsocketEventEnum.NewNotificationEvent]: z.object({
		messageId: z.string(),
		eventName: z.literal(WebsocketEventEnum.NewNotificationEvent),
		data: z.object({
			notificationId: z.string()
		})
	}),
	[WebsocketEventEnum.SystemReloadEvent]: z.object({
		messageId: z.string(),
		eventName: z.literal(WebsocketEventEnum.SystemReloadEvent),
		data: z.object({
			messageText: z.string(),
			messageTitle: z.string(),
			isReloadRequired: z.boolean()
		})
	}),
	[WebsocketEventEnum.ConversationAssignmentEvent]: z.object({
		messageId: z.string(),
		eventName: z.literal(WebsocketEventEnum.ConversationAssignmentEvent),
		data: z.object({
			assignedToMemberId: z.string(),
			conversationId: z.string(),
			assignedAt: z.string()
		})
	}),
	[WebsocketEventEnum.ConversationClosedEvent]: z.object({
		messageId: z.string(),
		eventName: z.literal(WebsocketEventEnum.ConversationClosedEvent),
		data: z.object({
			conversationId: z.string()
		})
	}),
	[WebsocketEventEnum.NewConversationEvent]: z.object({
		messageId: z.string(),
		eventName: z.literal(WebsocketEventEnum.NewConversationEvent),
		data: z.object({
			conversationId: z.string()
		})
	}),
	[WebsocketEventEnum.MessageAcknowledgementEvent]: z.object({
		eventName: z.literal(WebsocketEventEnum.MessageAcknowledgementEvent),
		messageId: z.string(),
		data: z.object({
			message: z.string()
		})
	}),
	[WebsocketEventEnum.PingEvent]: z.object({
		eventName: z.literal(WebsocketEventEnum.PingEvent),
		data: z.object({
			message: z.string()
		})
	})
}
