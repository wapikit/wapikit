import type { z } from 'zod'
import type { ApiServerEventDataMap, ApiServerEventEnum } from '../api-server-events'
import { type ConversationInboxStoreType } from '~/store/conversation-inbox.store'
import { type ConversationSchema } from 'root/.generated'

export function messageEventHandler(params: {
	conversations: ConversationSchema[]
	eventData: z.infer<(typeof ApiServerEventDataMap)[ApiServerEventEnum.NewMessageEvent]>
	writeProperty: ConversationInboxStoreType['writeProperty']
}): boolean {
	try {
		const { conversations, eventData, writeProperty } = params
		const conversation = conversations.find(
			convo => convo.uniqueId === eventData.message.conversationId
		)

		if (!conversation) {
			return false
		}

		const updatedConversation: ConversationSchema = {
			...conversation,
			// @ts-ignore - Will fix these types soon
			messages: [...conversation.messages, eventData.message]
		}

		writeProperty({
			conversations: conversations.map(convo =>
				convo.uniqueId === conversation.uniqueId ? updatedConversation : convo
			)
		})

		return true
	} catch (error) {
		console.error(error)
		return false
	}
}

export async function conversationAssignedEventHandler(
	message: z.infer<
		(typeof ApiServerEventDataMap)[ApiServerEventEnum.ConversationAssignmentEvent]['shape']['data']
	>
): Promise<boolean> {
	try {
		const { conversationId } = message
		console.log({ conversationId })

		// ! get the conversation from the store
		// ! update the conversation with the new assignee

		// ! show a notification that the conversation has been assigned if the conversation is assigned to the current user

		return true

		// ! append the above message to the conversation
	} catch (error) {
		console.error(error)
		return false
	}
}

export async function conversationUnassignedEventHandler(
	message: z.infer<
		(typeof ApiServerEventDataMap)[ApiServerEventEnum.ConversationClosedEvent]['shape']['data']
	>
) {
	try {
		const { conversationId } = message
		console.log({ conversationId })

		return true
	} catch (error) {
		console.error(error)
		return false
	}
}
