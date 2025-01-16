import { MessageDirectionEnum, MessageStatusEnum, MessageTypeEnum } from 'root/.generated'
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
		eventId: z.string(),
		data: z.object({
			conversationId: z.string(),
			createdAt: z.string(),
			direction: z.nativeEnum(MessageDirectionEnum),
			message_type: z.nativeEnum(MessageTypeEnum),
			messageData: z.record(z.string(), z.unknown()),
			status: z.nativeEnum(MessageStatusEnum),
			uniqueId: z.string()
		})
	}),
	[WebsocketEventEnum.NotificationReadEvent]: z.object({
		eventName: z.literal(WebsocketEventEnum.NotificationReadEvent),
		eventId: z.string(),
		data: z.object({
			userId: z.string(),
			organizationId: z.string(),
			notificationId: z.string()
		})
	}),
	[WebsocketEventEnum.MessageReadEvent]: z.object({
		eventId: z.string(),
		eventName: z.literal(WebsocketEventEnum.MessageReadEvent),
		data: z.object({
			messageId: z.string()
		})
	}),
	[WebsocketEventEnum.NewNotificationEvent]: z.object({
		eventId: z.string(),
		eventName: z.literal(WebsocketEventEnum.NewNotificationEvent),
		data: z.object({
			userId: z.string(),
			organizationId: z.string(),
			notificationPayload: z.object({
				createdAt: z.string(),
				description: z.string(),
				read: z.boolean(),
				title: z.string(),
				type: z.string(),
				uniqueId: z.string(),
				ctaUrl: z.string().optional()
			})
		})
	}),
	[WebsocketEventEnum.SystemReloadEvent]: z.object({
		eventId: z.string(),
		eventName: z.literal(WebsocketEventEnum.SystemReloadEvent),
		data: z.object({
			messageText: z.string(),
			messageTitle: z.string(),
			isReloadRequired: z.boolean()
		})
	}),
	[WebsocketEventEnum.ConversationAssignmentEvent]: z.object({
		eventId: z.string(),
		eventName: z.literal(WebsocketEventEnum.ConversationAssignmentEvent),
		data: z.object({
			assignedToMemberId: z.string(),
			conversationId: z.string(),
			assignedAt: z.string()
		})
	}),
	[WebsocketEventEnum.ConversationClosedEvent]: z.object({
		eventId: z.string(),
		eventName: z.literal(WebsocketEventEnum.ConversationClosedEvent),
		data: z.object({
			conversationId: z.string()
		})
	}),
	[WebsocketEventEnum.NewConversationEvent]: z.object({
		eventId: z.string(),
		eventName: z.literal(WebsocketEventEnum.NewConversationEvent),
		data: z.object({
			conversationId: z.string()
		})
	}),
	[WebsocketEventEnum.MessageAcknowledgementEvent]: z.object({
		eventName: z.literal(WebsocketEventEnum.MessageAcknowledgementEvent),
		eventId: z.string(),
		data: z.object({
			message: z.string()
		})
	}),
	[WebsocketEventEnum.PingEvent]: z.object({
		eventName: z.literal(WebsocketEventEnum.PingEvent),
		eventId: z.string(),
		data: z.object({
			message: z.string()
		})
	})
}
