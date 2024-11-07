import { produce } from 'immer'
import { create } from 'zustand'

export type FeatureFlagStoreType = {
	writeProperty: (
		updates:
			| WritePropertyParamType
			| ((state?: FeatureFlagStoreType | undefined) => FeatureFlagStoreType)
	) => void
	resetStore: () => void
	integrationFeatureFlags: {
		isSlackIntegrationEnabled: boolean
		isOpenAiIntegrationEnabled: boolean
		isCustomChatBoxIntegrationEnabled: boolean
	}
	systemFeatureFlags: {
		isRoleBasedAccessControlEnabled: boolean
		isQuickActionsKeywordEnabled: boolean
		isMultiOrganizationConfigurationEnabled: boolean
		isApiAccessEnabled: boolean
	}
}

type WritePropertyParamType = {
	[K in keyof FeatureFlagStoreType]?: FeatureFlagStoreType[K]
}

const useFeatureFlagStore = create<FeatureFlagStoreType>(set => ({
	writeProperty: updates => {
		if (typeof updates === 'object') {
			set(state => ({
				...state,
				...updates
			}))
		} else {
			set(state => produce<FeatureFlagStoreType>(state, updates))
		}
	},
	resetStore: () => {
		set(() => ({}))
	},
	integrationFeatureFlags: {
		isSlackIntegrationEnabled: false,
		isOpenAiIntegrationEnabled: false,
		isCustomChatBoxIntegrationEnabled: false
	},
	systemFeatureFlags: {
		isRoleBasedAccessControlEnabled: false,
		isQuickActionsKeywordEnabled: false,
		isMultiOrganizationConfigurationEnabled: false,
		isApiAccessEnabled: false
	}
}))

export { useFeatureFlagStore as useHireFormStore }
