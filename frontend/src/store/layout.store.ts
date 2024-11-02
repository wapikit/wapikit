import { produce } from 'immer'
import { type GetUserResponseSchema } from 'root/.generated'
import { create } from 'zustand'

export type LayoutStoreType = {
	notifications: string[]
	isOwner: boolean
	user: GetUserResponseSchema | null
	writeProperty: (
		updates: WritePropertyParamType | ((state?: LayoutStoreType | undefined) => LayoutStoreType)
	) => void
	resetStore: () => void
}

type WritePropertyParamType = {
	[K in keyof LayoutStoreType]?: LayoutStoreType[K]
}

const useLayoutStore = create<LayoutStoreType>(set => ({
	notifications: [],
	isOwner: false,
	user: null,
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

export { useLayoutStore }
