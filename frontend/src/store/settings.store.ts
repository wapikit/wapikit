import { produce } from 'immer'
import { OrganizationRoleSchema } from 'root/.generated'
import { create } from 'zustand'

export type SettingsStoreType = {
	applicationSettings: {
		name: string | null
		logo: string | null
		favicon: string | null
	}
	organizationSettings: {
		name: string | null
		logo: string | null
	}
	whatsappSettings: {
		defaultPhoneNumber: string | null
		phoneNumbers: {
			number: string
			isDefault: boolean
		}[]
	}
	quickReplies: {
		id: string
		message: string
		reply: string
	}[]
	apiKey: string
	roles: OrganizationRoleSchema[]
	writeProperty: (
		updates:
			| WritePropertyParamType
			| ((state?: SettingsStoreType | undefined) => SettingsStoreType)
	) => void
	resetStore: () => void
}

type WritePropertyParamType = {
	[K in keyof SettingsStoreType]?: SettingsStoreType[K]
}

const useSettingsStore = create<SettingsStoreType>(set => ({
	organizationSettings: {
		logo: '',
		name: ''
	},
	whatsappSettings: {
		defaultPhoneNumber: '',
		phoneNumbers: []
	},
	applicationSettings: {
		favicon: '',
		logo: '',
		name: ''
	},
	quickReplies: [],
	roles: [],
	apiKey: '',
	writeProperty: updates => {
		if (typeof updates === 'object') {
			set(state => ({
				...state,
				...updates
			}))
		} else {
			set(state => produce<SettingsStoreType>(state, updates))
		}
	},
	resetStore: () => {
		set(() => ({}))
	}
}))

export { useSettingsStore }
