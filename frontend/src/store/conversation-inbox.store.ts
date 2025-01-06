import { produce } from 'immer'
import { type ConversationSchema } from 'root/.generated'
import { create } from 'zustand'

export type ConversationInboxStoreType = {
	writeProperty: (
		updates:
			| WritePropertyParamType
			| ((state: ConversationInboxStoreType) => ConversationInboxStoreType)
	) => void
	resetStore: () => void
	conversations: ConversationSchema[]
	currentConversation: ConversationSchema | null
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
	currentConversation: null,
	conversations: []
}))

export { useConversationInboxStore }
