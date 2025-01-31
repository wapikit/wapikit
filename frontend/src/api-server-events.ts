
import { z } from 'zod'

export enum ApiServerEventEnum {
	MessageAcknowledgementEvent = 'MessageAcknowledgementEvent',
	NewMessageEvent = 'NewMessage',
	NotificationReadEvent = 'NotificationReadEvent',
	MessageReadEvent = 'MessageReadEvent',
	NewNotificationEvent = 'NewNotificationEvent',
	SystemReloadEvent = 'SystemReloadEvent',
	ConversationAssignmentEvent = 'ConversationAssignmentEvent',
	ConversationClosedEvent = 'ConversationClosedEvent',
	NewConversationEvent = 'NewConversationEvent',
	PingEvent = 'PingEvent'
}

export const ApiServerEventDataMap = {
	[ApiServerEventEnum.NewMessageEvent]: z.object({
		conversation: z.record(z.string(), z.unknown()),
		message: z.record(z.string(), z.unknown()),
	}),
	[ApiServerEventEnum.NotificationReadEvent]: z.object({
		eventName: z.literal(ApiServerEventEnum.NotificationReadEvent),
		data: z.object({
			userId: z.string(),
			organizationId: z.string(),
			notificationId: z.string()
		})
	}),
	[ApiServerEventEnum.MessageReadEvent]: z.object({

		eventName: z.literal(ApiServerEventEnum.MessageReadEvent),
		data: z.object({
			messageId: z.string()
		})
	}),
	[ApiServerEventEnum.NewNotificationEvent]: z.object({
		eventName: z.literal(ApiServerEventEnum.NewNotificationEvent),
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
	[ApiServerEventEnum.SystemReloadEvent]: z.object({

		eventName: z.literal(ApiServerEventEnum.SystemReloadEvent),
		data: z.object({
			messageText: z.string(),
			messageTitle: z.string(),
			isReloadRequired: z.boolean()
		})
	}),
	[ApiServerEventEnum.ConversationAssignmentEvent]: z.object({

		eventName: z.literal(ApiServerEventEnum.ConversationAssignmentEvent),
		data: z.object({
			assignedToMemberId: z.string(),
			conversationId: z.string(),
			assignedAt: z.string()
		})
	}),
	[ApiServerEventEnum.ConversationClosedEvent]: z.object({

		eventName: z.literal(ApiServerEventEnum.ConversationClosedEvent),
		data: z.object({
			conversationId: z.string()
		})
	}),
	[ApiServerEventEnum.NewConversationEvent]: z.object({

		eventName: z.literal(ApiServerEventEnum.NewConversationEvent),
		data: z.object({
			conversationId: z.string()
		})
	}),
	[ApiServerEventEnum.MessageAcknowledgementEvent]: z.object({
		eventName: z.literal(ApiServerEventEnum.MessageAcknowledgementEvent),

		data: z.object({
			message: z.string()
		})
	}),
	[ApiServerEventEnum.PingEvent]: z.object({
		eventName: z.literal(ApiServerEventEnum.PingEvent),

		data: z.object({
			message: z.string()
		})
	})
}