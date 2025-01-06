import type { z } from 'zod'
import type { WebsocketEventDataMap, WebsocketEventEnum } from '../websocket-events'
import { type ConversationInboxStoreType } from '~/store/conversation-inbox.store'
import { type ConversationSchema } from 'root/.generated'

export async function messageEventHandler(params: {
	conversation: ConversationSchema
	message: z.infer<
		(typeof WebsocketEventDataMap)[WebsocketEventEnum.MessageEvent]['shape']['data']
	>
	writeProperty: ConversationInboxStoreType['writeProperty']
}): Promise<boolean> {
	try {
		const { conversation, message, writeProperty } = params

		console.log({ conversation, message })

		const updatedConversation: ConversationSchema = {
			...conversation,
			messages: [
				...conversation.messages
				// message
			]
		}

		writeProperty(data => {
			return {
				...data,
				conversations: data.conversations.map(convo => {
					if (convo.uniqueId === conversation.uniqueId) {
						return updatedConversation
					}
					return convo
				})
			}
		})

		await Promise.resolve()

		// ! get the conversation from the store
		// ! append the above message to the conversation
		// ! return true if the message was appended successfully

		return true
	} catch (error) {
		console.error(error)
		return false
	}
}

export async function conversationAssignedEventHandler(
	message: z.infer<
		(typeof WebsocketEventDataMap)[WebsocketEventEnum.ConversationAssignmentEvent]['shape']['data']
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
		(typeof WebsocketEventDataMap)[WebsocketEventEnum.ConversationClosedEvent]['shape']['data']
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
