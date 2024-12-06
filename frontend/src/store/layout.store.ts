import { produce } from 'immer'
import {
	type UserSchema,
	type GetUserResponseSchema,
	type GetAllPhoneNumbersResponseSchema,
	type GetAllMessageTemplatesResponseSchema
} from 'root/.generated'
import { create } from 'zustand'
import { OnboardingSteps } from '~/constants'

export type LayoutStoreType = {
	onboardingSteps: typeof OnboardingSteps
	notifications: string[]
	isOwner: boolean
	user: Omit<UserSchema, 'organization'> | null
	currentOrganization: GetUserResponseSchema['user']['organization'] | null
	phoneNumbers: GetAllPhoneNumbersResponseSchema
	templates: GetAllMessageTemplatesResponseSchema
	writeProperty: (
		updates: WritePropertyParamType | ((state?: LayoutStoreType | undefined) => LayoutStoreType)
	) => void
	resetStore: () => void
}

type WritePropertyParamType = {
	[K in keyof LayoutStoreType]?: LayoutStoreType[K]
}

const useLayoutStore = create<LayoutStoreType>(set => ({
	onboardingSteps: OnboardingSteps,
	notifications: [],
	isOwner: false,
	user: null,
	currentOrganization: null,
	phoneNumbers: [],
	templates: [],
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
