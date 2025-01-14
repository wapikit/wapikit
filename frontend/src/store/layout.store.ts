import { produce } from 'immer'
import {
	type UserSchema,
	type GetUserResponseSchema,
	type GetAllPhoneNumbersResponseSchema,
	type GetAllMessageTemplatesResponseSchema,
	type ContactSchema
} from 'root/.generated'
import { create } from 'zustand'
import { OnboardingSteps } from '~/constants'

export type LayoutStoreType = {
	featureFlags: {
		integrations: {
			isSlackIntegrationEnabled: boolean
			isOpenAiIntegrationEnabled: boolean
			isCustomChatBoxIntegrationEnabled: boolean
		}
		system: {
			isRoleBasedAccessControlEnabled: boolean
			isMultiOrganizationConfigurationEnabled: boolean
		}
	}
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
	contactSheetData: null,
	onboardingSteps: OnboardingSteps,
	notifications: [],
	isOwner: false,
	user: null,
	currentOrganization: null,
	phoneNumbers: [],
	templates: [],
	featureFlags: {
		integrations: {
			isSlackIntegrationEnabled: false,
			isOpenAiIntegrationEnabled: false,
			isCustomChatBoxIntegrationEnabled: false
		},
		system: {
			isRoleBasedAccessControlEnabled: false,
			isMultiOrganizationConfigurationEnabled: false
		}
	},
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
