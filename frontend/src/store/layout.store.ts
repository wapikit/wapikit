import { produce } from 'immer'
import {
	type UserSchema,
	type GetUserResponseSchema,
	type GetAllPhoneNumbersResponseSchema,
	type GetAllMessageTemplatesResponseSchema,
	type ContactSchema,
	type GetFeatureFlagsResponseSchema
} from 'root/.generated'
import { create } from 'zustand'
import { OnboardingSteps } from '~/constants'

export type LayoutStoreType = {
	playNotificationSound: () => void
	featureFlags: GetFeatureFlagsResponseSchema['featureFlags'] | null
	onboardingSteps: typeof OnboardingSteps
	notifications: string[]
	isOwner: boolean
	user: Omit<UserSchema, 'organization'> | null
	currentOrganization: GetUserResponseSchema['user']['organization'] | null
	phoneNumbers: GetAllPhoneNumbersResponseSchema
	templates: GetAllMessageTemplatesResponseSchema
	contactSheetData: ContactSchema | null
	writeProperty: (
		updates: WritePropertyParamType | ((state?: LayoutStoreType | undefined) => LayoutStoreType)
	) => void
	resetStore: () => void
}

type WritePropertyParamType = {
	[K in keyof LayoutStoreType]?: LayoutStoreType[K]
}

const useLayoutStore = create<LayoutStoreType>(set => ({
	playNotificationSound() {
		const audio = new Audio('/assets/notification-sounds/pop.wav')
		audio.play().catch(error => console.error(error))
	},
	isAiChatBoxOpen: false,
	contactSheetData: null,
	onboardingSteps: OnboardingSteps,
	notifications: [],
	isOwner: false,
	user: null,
	currentOrganization: null,
	phoneNumbers: [],
	templates: [],
	featureFlags: null,
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
