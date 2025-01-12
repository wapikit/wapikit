import { produce } from 'immer'
import { create } from 'zustand'

export type AiChatBoxStoreType = {
	writeProperty: (
		updates: WritePropertyParamType | ((state: AiChatBoxStoreType) => AiChatBoxStoreType)
	) => void
	resetStore: () => void
	isOpen: boolean
}

type WritePropertyParamType = {
	[K in keyof AiChatBoxStoreType]?: AiChatBoxStoreType[K]
}

const useAiChatBoxStore = create<AiChatBoxStoreType>(set => ({
	writeProperty: updates => {
		console.log('updating store', updates)
		if (typeof updates === 'object') {
			set(state => ({
				...state,
				...updates
			}))
		} else {
			set(state => produce<AiChatBoxStoreType>(state, updates))
		}
	},
	resetStore: () => {
		set(() => ({}))
	},
	isOpen: false
}))

export { useAiChatBoxStore }
