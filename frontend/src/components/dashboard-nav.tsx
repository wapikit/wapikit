'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'

import { Icons } from '~/components/icons'
import { clsx as cn } from 'clsx'
import { type NavItem } from '~/types'
import { useState, type Dispatch, type SetStateAction } from 'react'
import { useSidebar } from '~/hooks/use-sidebar'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from './ui/tooltip'
import {
	useCreateOrganization,
	useGetUserOrganizations,
	useSwitchOrganization
} from 'root/.generated'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from './ui/select'
import { useAuthState } from '~/hooks/use-auth-state'
import { AUTH_TOKEN_LS } from '~/constants'
import { useForm } from 'react-hook-form'
import { type z } from 'zod'
import { NewOrganizationFormSchema } from '~/schema'
import { zodResolver } from '@hookform/resolvers/zod'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from './ui/form'
import { Modal } from './ui/modal'
import { errorNotification } from '~/reusable-functions'
import { Input } from './ui/input'
import { Button } from '~/components/ui/button'
import { Plus } from 'lucide-react'
import { useLayoutStore } from '~/store/layout.store'

interface DashboardNavProps {
	items: NavItem[]
	setOpen?: Dispatch<SetStateAction<boolean>>
	isMobileNav?: boolean
}

export function DashboardNav({ items, setOpen, isMobileNav = false }: DashboardNavProps) {
	const path = usePathname()
	const { isMinimized } = useSidebar()
	const { authState } = useAuthState()
	const { currentOrganization, writeProperty } = useLayoutStore()

	const [isNewOrganizationFormModalOpen, setIsNewOrganizationFormModalOpen] = useState(false)

	const newOrganizationForm = useForm<z.infer<typeof NewOrganizationFormSchema>>({
		resolver: zodResolver(NewOrganizationFormSchema),
		defaultValues: {
			name: '',
			description: ''
		}
	})

	const {
		isFetching,
		data: organizations,
		refetch: refetchOrganizations
	} = useGetUserOrganizations({
		page: 1,
		per_page: 50
	})

	const switchOrganizationMutation = useSwitchOrganization()
	const createOrganizationMutation = useCreateOrganization()

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
						window.localStorage.setItem(AUTH_TOKEN_LS, response.token)
						window.location.reload()
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

	async function handleCreateOrganization(data: z.infer<typeof NewOrganizationFormSchema>) {
		try {
			const response = await createOrganizationMutation.mutateAsync({
				data: {
					name: data.name
					// description: data.description || undefined
				}
			})

			if (response.organization) {
				setIsNewOrganizationFormModalOpen(false)
				await refetchOrganizations()
			} else {
				// show error message
				errorNotification({
					message: 'Organization creation failed'
				})
			}
		} catch (error) {
			console.error('error', error)
			errorNotification({
				message: 'Organization creation failed'
			})
		}
	}

	return (
		<nav className="grid items-start gap-2">
			<Modal
				title="Create New Organization"
				description=""
				isOpen={isNewOrganizationFormModalOpen}
				onClose={() => {
					setIsNewOrganizationFormModalOpen(false)
				}}
			>
				<div className="flex w-full items-center justify-end space-x-2 pt-6">
					<Form {...newOrganizationForm}>
						<form
							onSubmit={newOrganizationForm.handleSubmit(handleCreateOrganization)}
							className="w-full space-y-8"
						>
							<div className="flex flex-col gap-8">
								<FormField
									control={newOrganizationForm.control}
									name="name"
									render={({ field }) => (
										<FormItem>
											<FormLabel>Name</FormLabel>
											<FormControl>
												<Input
													placeholder="name"
													{...field}
													autoComplete="off"
												/>
											</FormControl>
											<FormMessage />
										</FormItem>
									)}
								/>

								<FormField
									control={newOrganizationForm.control}
									name="description"
									render={({ field }) => (
										<FormItem>
											<FormLabel>Description</FormLabel>
											<FormControl>
												<Input
													placeholder="Description (optional)"
													{...field}
													autoComplete="off"
												/>
											</FormControl>
											<FormMessage />
										</FormItem>
									)}
								/>
							</div>
							<Button className="ml-auto mr-0 w-full" type="submit">
								Create Organization
							</Button>
						</form>
					</Form>
				</div>
			</Modal>

			{/* organization selection dropdown here */}
			<div className={isMinimized ? 'hidden' : 'flex items-center justify-center'}>
				<Select
					disabled={isFetching}
					onValueChange={e => {
						switchOrganization(e).catch(error => console.error(error))
					}}
					value={currentOrganization?.uniqueId || 'no organizations'}
				>
					<SelectTrigger>
						<SelectValue placeholder="Select list" />
					</SelectTrigger>

					<SelectContent>
						{!organizations?.organizations ||
						organizations.organizations.length === 0 ? (
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
								<Button
									onClick={() => {
										setIsNewOrganizationFormModalOpen(true)
									}}
									variant={'secondary'}
									className="my-2 w-full"
								>
									<Plus className="size-5" /> Create New Organization
								</Button>
							</>
						)}
					</SelectContent>
				</Select>
			</div>

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
											'flex cursor-pointer items-center gap-2 overflow-hidden rounded-md py-2 text-sm font-medium hover:bg-accent hover:text-accent-foreground',
											path === item.href ? 'bg-accent' : 'transparent',
											item.disabled && 'cursor-not-allowed opacity-80'
										)}
										onClick={() => {
											if (setOpen) setOpen(false)
										}}
									>
										<Icon className={`ml-3 size-4`} />

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

			<Button
				className="ml-2 mt-2 flex w-[80%] gap-2 text-left"
				onClick={() => {
					writeProperty({
						isCommandMenuOpen: true
					})
				}}
			>
				Quick Action
				<div>⌘ K</div>
			</Button>
		</nav>
	)
}
