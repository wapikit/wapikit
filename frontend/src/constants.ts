import { type NavItem } from './types'

export const IS_PRODUCTION = process.env.NODE_ENV === 'production'
export const IS_DEVELOPMENT = process.env.NODE_ENV === 'development'

export const AUTH_TOKEN_LS = '__auth_token'

export const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:8000'

export const IMG_MAX_LIMIT = 10

export const navItems: NavItem[] = [
	{
		title: 'Dashboard',
		href: '/dashboard',
		icon: 'dashboard',
		label: 'Dashboard'
	},
	{
		title: 'Members',
		href: '/members',
		icon: 'laptop',
		label: 'Members'
	},
	{
		title: 'Chat',
		href: '/chat',
		icon: 'message',
		label: 'Chat'
	},
	{
		title: 'Media Library',
		href: '/media',
		icon: 'media',
		label: 'Media Library'
	},
	{
		title: 'Campaigns',
		href: '/campaigns',
		icon: 'rocket',
		label: 'Campaigns'
	},
	{
		title: 'Lists',
		href: '/lists',
		icon: 'billing',
		label: 'employee'
	},
	{
		title: 'Contacts',
		href: '/contacts',
		icon: 'user',
		label: 'profile'
	},
	{
		title: 'Settings',
		href: '/settings',
		icon: 'settings',
		label: 'kanban'
	}
]
