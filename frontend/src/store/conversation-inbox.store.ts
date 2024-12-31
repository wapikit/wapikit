import { produce } from 'immer'
import { create } from 'zustand'

export type ConversationInboxStoreType = {
	writeProperty: (
		updates:
			| WritePropertyParamType
			| ((state?: ConversationInboxStoreType | undefined) => ConversationInboxStoreType)
	) => void
	resetStore: () => void
	conversations: {
		unique_id: string
		name: string
		isOnline: boolean
		messages: string[]
		unreadMessages: number
	}[]
}

type WritePropertyParamType = {
	[K in keyof ConversationInboxStoreType]?: ConversationInboxStoreType[K]
}

const useConversationInboxStore = create<ConversationInboxStoreType>(set => ({
	writeProperty: updates => {
		if (typeof updates === 'object') {
			set(state => ({
				...state,
				...updates
			}))
		} else {
			set(state => produce<ConversationInboxStoreType>(state, updates))
		}
	},
	resetStore: () => {
		set(() => ({}))
	},
	activeConversation: {
		isOnline: true,
		messages: []
	},
	conversations: []
}))

export { useConversationInboxStore as useHireFormStore }
