import type { Icons } from '~/components/icons'

export interface NavItem {
	title: string
	href?: string
	disabled?: boolean
	external?: boolean
	icon?: keyof typeof Icons
	label?: string
	description?: string
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

export enum WebsocketStatusEnum {
	Connecting = 'connecting',
	Connected = 'connected',
	Disconnected = 'disconnected',
	Idle = 'idle'
}
