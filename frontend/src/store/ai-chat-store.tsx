import { produce } from 'immer'
import { type AiChatMessageSchema, type AiChatSchema } from 'root/.generated'
import { create } from 'zustand'

export type AiChatStoreType = {
	writeProperty: (
		updates: WritePropertyParamType | ((state: AiChatStoreType) => AiChatStoreType)
	) => void
	inputValue: string
	resetStore: () => void
	pushMessage(message: AiChatMessageSchema): void
	updateChatMessage: (messageId: string, content: string) => void
	updateUserMessageId: (messageId: string) => void
	editMessage: (messageId: string, content: string) => void
	isOpen: boolean
	chats: AiChatSchema[]
	currentChatMessages: AiChatMessageSchema[]
	suggestions: {
		title: string
		label: string
		action: string
	}[]
}

type WritePropertyParamType = {
	[K in keyof AiChatStoreType]?: AiChatStoreType[K]
}

const useAiChatStore = create<AiChatStoreType>(set => ({
	inputValue: '',
	writeProperty: updates => {
		if (typeof updates === 'object') {
			set(state => ({
				...state,
				...updates
			}))
		} else {
			set(state => produce<AiChatStoreType>(state, updates))
		}
	},
	resetStore: () => {
		set(() => ({}))
	},
	pushMessage(message: AiChatMessageSchema) {
		set(state => ({
			...state,
			currentChatMessages: [...state.currentChatMessages, message]
		}))
	},
	updateChatMessage(messageId: string, content: string) {
		set(state => {
			const message = state.currentChatMessages.find(
				message => message.uniqueId === messageId
			)
			if (!message) return state

			const updatedMessage = {
				...message,
				content: message.content + content
			}

			const updatedMessages = state.currentChatMessages.map(message =>
				message.uniqueId === messageId ? updatedMessage : message
			)

			return {
				...state,
				currentChatMessages: updatedMessages
			}
		})
	},
	editMessage(messageId: string, content: string) {
		set(state => {
			const message = state.currentChatMessages.find(
				message => message.uniqueId === messageId
			)

			const messageIndex = state.currentChatMessages.findIndex(
				message => message.uniqueId === messageId
			)

			const isLastMessage = messageIndex === state.currentChatMessages.length - 1

			if (!message) return state

			const updatedMessage = {
				...message,
				content
			}

			if (isLastMessage) {
				return {
					...state,
					currentChatMessages: state.currentChatMessages.map(message =>
						message.uniqueId === messageId ? updatedMessage : message
					)
				}
			} else {
				// * if not last index then remove all the follow up messages and delete all the trailing messages
				const updatedMessages = state.currentChatMessages
					.map(message => (message.uniqueId === messageId ? updatedMessage : message))
					.filter((_, index) => index <= messageIndex)

				return {
					...state,
					currentChatMessages: updatedMessages
				}
			}
		})
	},
	updateUserMessageId(messageId: string) {
		// update the id of the last message

		set(state => {
			const lastMessage = state.currentChatMessages[state.currentChatMessages.length - 1]

			if (!lastMessage) return state

			const updatedMessage = {
				...lastMessage,
				uniqueId: messageId
			}

			const updatedMessages = state.currentChatMessages.map(message =>
				message.uniqueId === lastMessage.uniqueId ? updatedMessage : message
			)

			return {
				...state,
				currentChatMessages: updatedMessages
			}
		})
	},
	currentChatMessages: [],
	isOpen: false,
	chats: [],
	suggestions: [
		{
			title: 'Give me insights of ',
			label: 'last weeks campaigns.',
			action: 'Give me the insights of last weeks campaigns.'
		},
		{
			title: 'What do you think we should improve on',
			label: `to increase our open-rate?`,
			action: `What do you think we should improve on to increase our open-rate?`
		}
	]
}))

export { useAiChatStore }
