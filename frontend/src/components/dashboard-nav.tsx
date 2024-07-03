'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'

import { Icons } from '~/components/icons'
import { clsx as cn } from 'clsx'
import { type NavItem } from '~/types'
import { type Dispatch, type SetStateAction } from 'react'
import { useSidebar } from '~/hooks/use-sidebar'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from './ui/tooltip'
import { useGetUserOrganizations, useSwitchOrganization } from 'root/.generated'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from './ui/select'
import { useAuthState } from '~/hooks/use-auth-state'
import { useLocalStorage } from '~/hooks/use-local-storage'
import { AUTH_TOKEN_LS } from '~/constants'
import { useRouter } from 'next/navigation'

interface DashboardNavProps {
	items: NavItem[]
	setOpen?: Dispatch<SetStateAction<boolean>>
	isMobileNav?: boolean
}

export function DashboardNav({ items, setOpen, isMobileNav = false }: DashboardNavProps) {
	const path = usePathname()
	const { isMinimized } = useSidebar()
	const { authState } = useAuthState()
	const setLocalStorageState = useLocalStorage<string>(AUTH_TOKEN_LS, '')[1]
	const router = useRouter()

	const { isFetching, data: organizations } = useGetUserOrganizations({
		page: 1,
		per_page: 50
	})

	const switchOrganizationMutation = useSwitchOrganization()

	if (!items?.length) {
		return null
	}

	async function switchOrganization(selectedOrganizationId: string) {
		try {
			if (authState.isAuthenticated) {
				if (authState.data.user.organizationId === selectedOrganizationId) {
					// same selector clicked
					return
				} else {
					const response = await switchOrganizationMutation.mutateAsync({
						data: {
							organizationId: selectedOrganizationId
						}
					})

					if (response.token) {
						setLocalStorageState(response.token)
						router.refresh()
					} else {
						// error show error message sonner
					}
				}
			} else {
				console.log('not authenticated')
			}
		} catch (error) {
			console.error('error', error)
		}
	}

	return (
		<nav className="grid items-start gap-2">
			{/* organization selection dropdown here */}

			<Select
				disabled={isFetching}
				onValueChange={e => {
					switchOrganization(e).catch(error => console.error(error))
				}}
				value={organizations?.organizations?.[0]?.uniqueId || 'no organizations'}
			>
				<SelectTrigger>
					<SelectValue placeholder="Select list" />
				</SelectTrigger>

				<SelectContent>
					{!organizations?.organizations || organizations.organizations.length === 0 ? (
						<SelectItem value={'no list'} disabled>
							No organizations created yet.
						</SelectItem>
					) : (
						<>
							{organizations.organizations.map(org => (
								<SelectItem key={org.uniqueId} value={org.uniqueId}>
									{org.name}
								</SelectItem>
							))}
						</>
					)}
				</SelectContent>
			</Select>

			<TooltipProvider>
				{items.map((item, index) => {
					const Icon = Icons[item.icon || 'arrowRight']

					return (
						item.href && (
							<Tooltip key={index}>
								<TooltipTrigger asChild>
									<Link
										href={item.disabled ? '/' : item.href}
										className={cn(
											'flex items-center gap-2 overflow-hidden rounded-md py-2 text-sm font-medium hover:bg-accent hover:text-accent-foreground',
											path === item.href ? 'bg-accent' : 'transparent',
											item.disabled && 'cursor-not-allowed opacity-80'
										)}
										onClick={() => {
											if (setOpen) setOpen(false)
										}}
									>
										<Icon className={`ml-3 size-5`} />

										{isMobileNav || (!isMinimized && !isMobileNav) ? (
											<span className="mr-2 truncate">{item.title}</span>
										) : (
											''
										)}
									</Link>
								</TooltipTrigger>
								<TooltipContent
									align="center"
									side="right"
									sideOffset={8}
									className={!isMinimized ? 'hidden' : 'inline-block'}
								>
									{item.title}
								</TooltipContent>
							</Tooltip>
						)
					)
				})}
			</TooltipProvider>
		</nav>
	)
}
