import { produce } from 'immer'
import { create } from 'zustand'

export type LayoutStoreType = {
    notifications: string[]
    writeProperty: (
        updates:
            | WritePropertyParamType
            | ((state?: LayoutStoreType | undefined) => LayoutStoreType)
    ) => void
    resetStore: () => void
}

type WritePropertyParamType = {
    [K in keyof LayoutStoreType]?: LayoutStoreType[K]
}

const useHireFormStore = create<LayoutStoreType>(set => ({
    notifications: [],
    writeProperty: updates => {
        if (typeof updates === 'object') {
            set(state => ({
                ...state,
                ...updates
            }))
        } else {
            set(state => produce<LayoutStoreType>(state, updates))
        }
    },
    resetStore: () => {
        set(() => ({}))
    }
}))

export { useHireFormStore }
