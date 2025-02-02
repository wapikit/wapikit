import type { Icons } from '~/components/icons'

export interface NavItem {
	title: string
	href?: string
	disabled?: boolean
	external?: boolean
	icon?: keyof typeof Icons
	label?: string
	description?: string,
	requiredFeatureFlag?: string[]
}

export interface Contact {
	name: string
	phone: string
	list: string[]
}

export interface TableCellActionProps {
	icon: keyof typeof Icons
	label: string
	onClick: (data: any) => Promise<void> | void
	disabled?: boolean
}

export enum SseEventSourceStateEnum {
	Connecting = 'Connecting',
	Connected = 'Connected',
	Disconnected = 'Disconnected'
}

export enum ChatBotStateEnum {
	Idle = 'Idle',
	Streaming = 'Streaming',
	Thinking = 'Thinking'
}
