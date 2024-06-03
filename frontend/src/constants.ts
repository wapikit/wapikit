import { type NavItem } from './types'

export const IS_PRODUCTION = process.env.NODE_ENV === 'production'
export const IS_DEVELOPMENT = process.env.NODE_ENV === 'development'

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
        title: 'Campaigns',
        href: '/dashboard/campaigns',
        icon: 'user',
        label: 'user'
    },
    {
        title: 'Lists',
        href: '/dashboard/lists',
        icon: 'employee',
        label: 'employee'
    },
    {
        title: 'Contacts',
        href: '/dashboard/contacts',
        icon: 'profile',
        label: 'profile'
    },
    {
        title: 'Settings',
        href: '/dashboard/settings',
        icon: 'kanban',
        label: 'kanban'
    },
    {
        title: 'Members',
        href: '/dashboard/members',
        icon: 'login',
        label: 'login'
    }
]
