import { produce } from 'immer'
import { create } from 'zustand'

export type ConversationStoreType = {
    writeProperty: (
        updates:
            | WritePropertyParamType
            | ((state?: ConversationStoreType | undefined) => ConversationStoreType)
    ) => void
    resetStore: () => void
}

type WritePropertyParamType = {
    [K in keyof ConversationStoreType]?: ConversationStoreType[K]
}

const useHireFormStore = create<ConversationStoreType>(set => ({
    writeProperty: updates => {
        if (typeof updates === 'object') {
            set(state => ({
                ...state,
                ...updates
            }))
        } else {
            set(state => produce<ConversationStoreType>(state, updates))
        }
    },
    resetStore: () => {
        set(() => ({}))
    }
}))

export { useHireFormStore }
