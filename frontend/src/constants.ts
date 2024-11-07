import { type NavItem } from './types'

export const IS_PRODUCTION = process.env.NODE_ENV === 'production'
export const IS_DEVELOPMENT = process.env.NODE_ENV === 'development'

export const AUTH_TOKEN_LS = '__auth_token'

export const BACKEND_URL = process.env.BACKEND_URL || 'http://127.0.0.1:5000/api'

export const IMG_MAX_LIMIT = 10

export const navItems: NavItem[] = [
	{
		title: 'Dashboard',
		href: '/dashboard',
		icon: 'dashboard',
		label: 'Dashboard'
	},
	{
		title: 'Contacts',
		href: '/contacts',
		icon: 'user',
		label: 'profile'
	},
	{
		title: 'Lists',
		href: '/lists',
		icon: 'billing',
		label: 'employee'
	},
	{
		title: 'Members',
		href: '/members',
		icon: 'laptop',
		label: 'Members'
	},
	// {
	// 	title: 'Chat',
	// 	href: '/chat',
	// 	icon: 'message',
	// 	label: 'Chat'
	// },
	{
		title: 'Campaigns',
		href: '/campaigns',
		icon: 'rocket',
		label: 'Campaigns'
	},
	{
		title: 'Settings',
		href: '/settings',
		icon: 'settings',
		label: 'kanban'
	},
	{
		title: 'Integrations',
		href: '/integrations',
		icon: 'link',
		label: 'Integrations'
	}
]
