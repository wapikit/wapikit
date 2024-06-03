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
	name: string,
	phone: string
	list: string[]
}
