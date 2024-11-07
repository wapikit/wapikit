import { produce } from 'immer'
import { create } from 'zustand'

export type ChatInboxStoreType = {
	writeProperty: (
		updates:
			| WritePropertyParamType
			| ((state?: ChatInboxStoreType | undefined) => ChatInboxStoreType)
	) => void
	resetStore: () => void
	conversations: {
		unique_id: string
		name: string
		isOnline: boolean
		messages: string[]
		unreadMessages: number
	}[]
	activeConversation: {
		isOnline: boolean
		messages: string[]
	}
}

type WritePropertyParamType = {
	[K in keyof ChatInboxStoreType]?: ChatInboxStoreType[K]
}

const useHireFormStore = create<ChatInboxStoreType>(set => ({
	writeProperty: updates => {
		if (typeof updates === 'object') {
			set(state => ({
				...state,
				...updates
			}))
		} else {
			set(state => produce<ChatInboxStoreType>(state, updates))
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

export { useHireFormStore }
