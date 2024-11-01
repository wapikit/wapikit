'use client'

import { Button } from '~/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '~/components/ui/card'
import { ScrollArea } from '~/components/ui/scroll-area'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '~/components/ui/tabs'
import { Input } from '~/components/ui/input'
import RolesTable from '~/components/settings/roles-table'
import { useEffect, useRef, useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { useLayoutStore } from '~/store/layout.store'
import { useSettingsStore } from '~/store/settings.store'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '~/components/ui/tooltip'
import { errorNotification, materialConfirm, successNotification } from '~/reusable-functions'
import Image from 'next/image'
import { FileUploaderComponent } from '~/components/file-uploader'
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue
} from '~/components/ui/select'
import { Heading } from '~/components/ui/heading'
import { Separator } from '~/components/ui/separator'
import { Plus } from 'lucide-react'
import {
	RolePermissionEnum,
	useCreateOrganizationRole,
	useGetOrganizationRoleById,
	useUpdateOrganizationRoleById
} from 'root/.generated'
import { Modal } from '~/components/ui/modal'
import { NewRoleFormSchema } from '~/schema'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { type z } from 'zod'
import {
	Form,
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage
} from '~/components/ui/form'
import { MultiSelect } from '~/components/multi-select'
import { listStringEnumMembers } from 'ts-enum-utils'

export default function SettingsPage() {
	const tabs = [
		{
			slug: 'account',
			title: 'Account'
		},
		{
			slug: 'organization',
			title: 'Organization'
		},
		{
			slug: 'app-settings',
			title: 'App Settings'
		},
		{
			slug: 'whatsapp-business-account',
			title: 'WhatsApp Settings'
		},
		// {
		// 	slug: 'quick-actions',
		// 	title: 'Quick Actions'
		// },
		{
			slug: 'api-keys',
			title: 'API Keys'
		},
		{
			slug: 'rbac',
			title: 'Access Control (RBAC)'
		}
	]

	const searchParams = useSearchParams()
	const router = useRouter()
	const rolesDataSetRef = useRef(false)

	const [isRoleCreationModelOpen, setIsRoleCreationModelOpen] = useState(false)
	const [roleIdToEdit, setRoleIdToEdit] = useState<string | null>(null)
	const [activeTab, setActiveTab] = useState(
		searchParams.get('tab')?.toString() || 'app-settings'
	)

	const { isOwner } = useLayoutStore()
	const { organizationSettings, whatsappSettings } = useSettingsStore()
	// const updateOrganizationSettings = useUpdateSettings()
	const createRoleMutation = useCreateOrganizationRole()
	const updateRoleMutation = useUpdateOrganizationRoleById()

	const { data: roleData } = useGetOrganizationRoleById('', {
		query: {
			enabled: !!roleIdToEdit
		}
	})

	const form = useForm<z.infer<typeof NewRoleFormSchema>>({
		resolver: zodResolver(NewRoleFormSchema),
		defaultValues: roleData
			? {
					name: roleData.role.name,
					description: roleData.role.description,
					permissions: roleData.role.permissions
				}
			: {
					name: '',
					description: '',
					permissions: []
				}
	})

	useEffect(() => {
		if (rolesDataSetRef.current) return
		if (roleData) {
			form.reset({
				name: roleData.role.name,
				description: roleData.role.description,
				permissions: roleData.role.permissions
			})
			rolesDataSetRef.current = true
		}
	}, [roleData, form])

	useEffect(() => {
		if (searchParams.get('tab')) {
			setActiveTab(searchParams.get('tab')?.toString() || 'account')
		}
	}, [searchParams])

	useEffect(() => {
		if (roleIdToEdit) {
			setIsRoleCreationModelOpen(true)
		}
	}, [roleIdToEdit])

	// async function handleSettingsUpdate() {
	// 	await updateOrganizationSettings.mutateAsync({
	// 		data: {}
	// 	})
	// }

	const [isBusy, setIsBusy] = useState(false)

	async function deleteOrganization() {
		try {
			const confirmed = await materialConfirm({
				title: 'Delete Organization',
				description: 'Are you sure you want to delete this organization?'
			})

			if (!confirmed) {
				return
				// delete organization
			}
		} catch {
			console.error('Error deleting organization')
			errorNotification({
				message: 'Error deleting organization'
			})
		}
	}

	async function leaveOrganization() {
		try {
			const confirmed = await materialConfirm({
				title: 'Leave Organization',
				description: 'Are you sure you want to leave this organization?'
			})

			if (!confirmed) {
				return
				// delete organization
			}
		} catch {
			console.error('Error leaving organization')
			errorNotification({
				message: 'Error leaving organization'
			})
		}
	}

	async function submitRoleForm(data: z.infer<typeof NewRoleFormSchema>) {
		console.log('submit role form called with data', data)
		try {
			setIsBusy(true)
			//  ! check here if the role is being updated or created
			if (roleIdToEdit) {
				const response = await updateRoleMutation.mutateAsync({
					id: roleIdToEdit,
					data: {
						name: data.name,
						permissions: data.permissions,
						description: data.description || undefined
					}
				})

				if (response) {
					successNotification({
						message: 'Role updated successfully'
					})
					setIsRoleCreationModelOpen(false)
				} else {
					errorNotification({
						message: 'Error updating role'
					})
				}
			} else {
				const response = await createRoleMutation.mutateAsync({
					data: {
						name: data.name,
						permissions: data.permissions,
						description: data.description || undefined
					}
				})

				if (response) {
					successNotification({
						message: 'Role created successfully'
					})
					setIsRoleCreationModelOpen(false)
					setRoleIdToEdit(null)
				} else {
					errorNotification({
						message: 'Error creating role'
					})
				}
			}
		} catch (error) {
			console.error(error)
			errorNotification({
				message: 'Error creating / updating role'
			})
		} finally {
			setIsBusy(false)
		}
	}

	return (
		<ScrollArea className="h-full pr-8">
			<div className="flex-1 space-y-4 p-4 pt-6 md:p-8">
				<div className="flex items-center justify-between space-y-2">
					<h2 className="text-3xl font-bold tracking-tight">Settings</h2>
				</div>
				<Tabs defaultValue={activeTab} className="space-y-4">
					<TabsList>
						{tabs.map(tab => {
							return (
								<TabsTrigger
									key={tab.slug}
									value={tab.slug}
									onClick={() => {
										router.push(`/settings?tab=${tab.slug}`)
									}}
								>
									{tab.title}
								</TabsTrigger>
							)
						})}
					</TabsList>
					{tabs.map(tab => {
						return (
							<TabsContent
								key={tab.slug}
								value={tab.slug}
								className="space-y-4 py-10"
							>
								{tab.slug === 'app-settings' ? (
									<div className="mr-auto flex max-w-4xl flex-col gap-5">
										<Card>
											<CardHeader>
												<CardTitle>Application Name</CardTitle>
												<CardDescription>
													Used to identify your project in the dashboard.
												</CardDescription>
											</CardHeader>
											<CardContent>
												<form>
													<Input placeholder="Project Name" />
												</form>
											</CardContent>
										</Card>
										<Card className="flex flex-row">
											<div className="flex-1">
												<CardHeader>
													<CardTitle>Root Url</CardTitle>
													<CardDescription>
														Used to identify your project in the
														dashboard.
													</CardDescription>
												</CardHeader>
												<CardContent>
													<form>
														<Input placeholder="Project Name" />
													</form>
												</CardContent>
											</div>
											<div className="tremor-Divider-root mx-auto my-6 flex items-center justify-between gap-3 text-tremor-default text-tremor-content dark:text-dark-tremor-content">
												<div className="h-full w-[1px] bg-tremor-border dark:bg-dark-tremor-border"></div>
											</div>
											<div className="flex-1">
												<CardHeader>
													<CardTitle>Favicon Url </CardTitle>
													<CardDescription>
														Used to identify your project in the
														dashboard.
													</CardDescription>
												</CardHeader>
												<CardContent>
													<form>
														<Input placeholder="Project Name" />
													</form>
												</CardContent>
											</div>
										</Card>
										<Card className="flex flex-row">
											<div className="flex-1">
												<CardHeader>
													<CardTitle>Media Upload Path</CardTitle>
													<CardDescription>
														Used to identify your project in the
														dashboard.
													</CardDescription>
												</CardHeader>
												<CardContent>
													<form>
														<Input placeholder="Project Name" />
													</form>
												</CardContent>
											</div>
											<div className="tremor-Divider-root mx-auto my-6 flex items-center justify-between gap-3 text-tremor-default text-tremor-content dark:text-dark-tremor-content">
												<div className="h-full w-[1px] bg-tremor-border dark:bg-dark-tremor-border"></div>
											</div>
											<div className="flex-1">
												<CardHeader>
													<CardTitle>Media Upload URI</CardTitle>
													<CardDescription>
														Used to identify your project in the
														dashboard.
													</CardDescription>
												</CardHeader>
												<CardContent>
													<form>
														<Input placeholder="Project Name" />
													</form>
												</CardContent>
											</div>
										</Card>
									</div>
								) : tab.slug === 'whatsapp-business-account' ? (
									<div className="mr-auto flex max-w-4xl flex-col gap-5">
										<div className="flex flex-row gap-5">
											<Card className="flex flex-1 items-center justify-between">
												<CardHeader>
													<CardTitle>Sync Templates</CardTitle>
													<CardDescription>
														Click the button to sync your templates with
														WhatsApp Business Account.
													</CardDescription>
												</CardHeader>
												<CardContent className="flex h-fit items-center justify-center pb-0">
													<TooltipProvider>
														<Tooltip>
															<TooltipTrigger asChild>
																<Button
																	disabled={!isOwner}
																	onClick={() => {}}
																>
																	Sync
																</Button>
															</TooltipTrigger>
															<TooltipContent
																align="center"
																side="right"
																sideOffset={8}
																className={
																	!isOwner
																		? 'hidden'
																		: 'inline-block'
																}
															>
																You are not the owner of this
																organization.
															</TooltipContent>
														</Tooltip>
													</TooltipProvider>
												</CardContent>
											</Card>

											<Card className="flex flex-1 items-center justify-between">
												<CardHeader>
													<CardTitle>Sync Phone Numbers</CardTitle>
													<CardDescription>
														Click the button to sync your phone number
														with WhatsApp Business Account.
													</CardDescription>
												</CardHeader>
												<CardContent className="flex h-fit items-center justify-center pb-0">
													<TooltipProvider>
														<Tooltip>
															<TooltipTrigger asChild>
																<Button
																	onClick={() => {}}
																	disabled={!isOwner}
																>
																	Sync
																</Button>
															</TooltipTrigger>
															<TooltipContent
																align="center"
																side="right"
																sideOffset={8}
																className={
																	isOwner
																		? 'hidden'
																		: 'inline-block'
																}
															>
																You are not the owner of this
																organization.
															</TooltipContent>
														</Tooltip>
													</TooltipProvider>
												</CardContent>
											</Card>
										</div>

										<Card className="flex flex-1 items-center justify-between">
											<CardHeader>
												<CardTitle>Default Phone Number</CardTitle>
											</CardHeader>
											<CardContent className="flex h-fit items-center justify-center pb-0">
												<Select
													onValueChange={e => {
														console.log(e)
													}}
													value={
														whatsappSettings.defaultPhoneNumber ||
														'no organizations'
													}
												>
													<SelectTrigger>
														<SelectValue placeholder="Select Phone Number" />
													</SelectTrigger>

													<SelectContent>
														{!whatsappSettings.phoneNumbers ||
														whatsappSettings.phoneNumbers.length ===
															0 ? (
															<SelectItem value={'empty'} disabled>
																No Phone Numbers.
															</SelectItem>
														) : (
															<>
																{whatsappSettings.phoneNumbers.map(
																	phoneNumber => (
																		<SelectItem
																			key={phoneNumber.number}
																			value={
																				phoneNumber.number
																			}
																		>
																			{phoneNumber.number}
																		</SelectItem>
																	)
																)}
															</>
														)}
													</SelectContent>
												</Select>
											</CardContent>
										</Card>
									</div>
								) : tab.slug === 'account' ? (
									<div className="mr-auto flex max-w-4xl flex-col gap-5">
										<Card>
											<CardHeader>
												<CardTitle>Profile Picture</CardTitle>
											</CardHeader>
											<CardContent className="flex h-fit w-full items-center justify-center pb-0">
												<Image
													src={
														'https://www.creatorlens.co/assets/empty-pfp.png'
													}
													width={500}
													height={500}
													alt="profile"
													className="h-40 w-40 rounded-full"
												/>
												<div className="flex-1">
													<FileUploaderComponent
														descriptionString="JPG / JPEG / PNG"
														onFileUpload={() => {
															console.log('file uploaded')
														}}
													/>
												</div>
											</CardContent>
										</Card>
										<Card className="flex flex-row">
											<div className="flex-1">
												<CardHeader>
													<CardTitle>Name</CardTitle>
												</CardHeader>
												<CardContent>
													<form>
														<Input placeholder="Name" />
													</form>
												</CardContent>
											</div>
											<div className="tremor-Divider-root mx-auto my-6 flex items-center justify-between gap-3 text-tremor-default text-tremor-content dark:text-dark-tremor-content">
												<div className="h-full w-[1px] bg-tremor-border dark:bg-dark-tremor-border"></div>
											</div>
											<div className="flex-1">
												<CardHeader>
													<CardTitle>Email</CardTitle>
												</CardHeader>
												<CardContent>
													<form>
														<Input placeholder="Email" type="email" />
													</form>
												</CardContent>
											</div>
										</Card>
										<Card className="flex flex-1 items-center justify-between">
											<CardHeader>
												<CardTitle>Delete Account</CardTitle>
											</CardHeader>
											<CardContent className="flex h-fit items-center justify-center pb-0">
												<TooltipProvider>
													<Tooltip>
														<TooltipTrigger asChild>
															<Button
																variant={'destructive'}
																onClick={() => {}}
															>
																Delete Account
															</Button>
														</TooltipTrigger>
														<TooltipContent
															align="center"
															side="right"
															sideOffset={8}
															className={
																!isOwner ? 'hidden' : 'inline-block'
															}
														>
															You are the owner of this organization.
														</TooltipContent>
													</Tooltip>
												</TooltipProvider>
											</CardContent>
										</Card>
									</div>
								) : tab.slug === 'quick-actions' ? (
									<></>
								) : tab.slug === 'api-keys' ? (
									<div className="mr-auto flex max-w-4xl flex-col gap-5">
										<Card className="border-none">
											<CardHeader>
												<CardTitle>API Access Key</CardTitle>
												<CardDescription>
													Use this API key to authenticate wapikit API
													requests.
												</CardDescription>
											</CardHeader>
											<CardContent>
												<form className="w-full max-w-sm">
													<Input
														placeholder="***********************"
														className="w-fit px-6"
														type="password"
														disabled
													/>
												</form>
											</CardContent>
										</Card>
									</div>
								) : tab.slug === 'rbac' ? (
									<div className="flex-1 space-y-4">
										<Modal
											title={roleIdToEdit ? 'Edit Role' : 'Create New Role'}
											description={
												roleIdToEdit
													? 'Edit the role for your organization'
													: 'Create a new role for your organization'
											}
											isOpen={isRoleCreationModelOpen}
											onClose={() => {
												setIsRoleCreationModelOpen(false)
												form.reset()
											}}
										>
											<div className="flex w-full items-center justify-end space-x-2 pt-6">
												<Form {...form}>
													<form
														onSubmit={form.handleSubmit(submitRoleForm)}
														className="w-full space-y-8"
													>
														<div className="flex flex-col gap-8">
															<FormField
																control={form.control}
																name="name"
																render={({ field }) => (
																	<FormItem>
																		<FormLabel>
																			Role Name
																		</FormLabel>
																		<FormControl>
																			<Input
																				disabled={isBusy}
																				placeholder="role name"
																				{...field}
																				autoComplete="off"
																			/>
																		</FormControl>
																		<FormMessage />
																	</FormItem>
																)}
															/>

															<FormField
																control={form.control}
																name="description"
																render={({ field }) => (
																	<FormItem>
																		<FormLabel>
																			Role Description
																		</FormLabel>
																		<FormControl>
																			<Input
																				disabled={isBusy}
																				placeholder="role description"
																				{...field}
																				autoComplete="off"
																			/>
																		</FormControl>
																		<FormMessage />
																	</FormItem>
																)}
															/>

															<FormField
																control={form.control}
																name="permissions"
																render={({}) => (
																	<FormItem className="tablet:w-3/4 tablet:gap-2 desktop:w-1/2 flex flex-col gap-1 ">
																		<FormLabel>
																			Select the permissions
																		</FormLabel>
																		<MultiSelect
																			options={listStringEnumMembers(
																				RolePermissionEnum
																			).map(item => {
																				return {
																					label: item.name,
																					value: item.value
																				}
																			})}
																			onValueChange={e => {
																				console.log({ e })
																				form.setValue(
																					'permissions',
																					e as RolePermissionEnum[],
																					{
																						shouldValidate:
																							true
																					}
																				)
																			}}
																			defaultValue={form.watch(
																				'permissions'
																			)}
																			placeholder="Select lists"
																			variant="default"
																		/>
																		<FormMessage />
																	</FormItem>
																)}
															/>
														</div>
														<Button
															disabled={isBusy}
															className="ml-auto mr-0 w-full"
															type="submit"
														>
															Create
														</Button>
													</form>
												</Form>
											</div>
										</Modal>

										<div className="flex items-start justify-between">
											<Heading
												title={`Manage Organization Roles`}
												description=""
											/>
											<div className="flex gap-2">
												<Button
													onClick={() => {
														// open the roles create modal
														setIsRoleCreationModelOpen(true)
													}}
												>
													<Plus className="mr-2 h-4 w-4" /> Add New
												</Button>
											</div>
										</div>
										<Separator />
										<RolesTable setRoleToEditId={setRoleIdToEdit} />
									</div>
								) : tab.slug === 'organization' ? (
									<div className="mr-auto flex max-w-4xl flex-col gap-5">
										{/* organization name update button */}
										<Card>
											<CardHeader>
												<CardTitle>Organization Name</CardTitle>
											</CardHeader>
											<CardContent>
												<form>
													<Input
														placeholder={
															organizationSettings.name ||
															'Organization Name'
														}
														disabled={!isOwner}
													/>
												</form>
											</CardContent>
										</Card>

										{/* leave organization button */}
										<div className="flex flex-row gap-5">
											<Card className="flex flex-1 items-center justify-between">
												<CardHeader>
													<CardTitle>Leave Organization</CardTitle>
												</CardHeader>
												<CardContent className="flex h-fit items-center justify-center pb-0">
													<TooltipProvider>
														<Tooltip>
															<TooltipTrigger asChild>
																<Button
																	variant={'destructive'}
																	disabled={isOwner}
																	onClick={() => {
																		leaveOrganization().catch(
																			error =>
																				console.error(error)
																		)
																	}}
																>
																	Leave
																</Button>
															</TooltipTrigger>
															<TooltipContent
																align="center"
																side="right"
																sideOffset={8}
																className={
																	!isOwner
																		? 'hidden'
																		: 'inline-block'
																}
															>
																You are the owner of this
																organization.
															</TooltipContent>
														</Tooltip>
													</TooltipProvider>
												</CardContent>
											</Card>

											{/* delete organization button */}
											<Card className="flex flex-1 items-center justify-between">
												<CardHeader>
													<CardTitle>Delete Organization</CardTitle>
												</CardHeader>
												<CardContent className="flex h-fit items-center justify-center pb-0">
													<TooltipProvider>
														<Tooltip>
															<TooltipTrigger asChild>
																<Button
																	variant={'destructive'}
																	onClick={() => {
																		deleteOrganization().catch(
																			error =>
																				console.error(error)
																		)
																	}}
																>
																	Delete
																</Button>
															</TooltipTrigger>
															<TooltipContent
																align="center"
																side="right"
																sideOffset={8}
																className={
																	isOwner
																		? 'hidden'
																		: 'inline-block'
																}
															>
																You are not the owner of this
																organization.
															</TooltipContent>
														</Tooltip>
													</TooltipProvider>
												</CardContent>
											</Card>
										</div>
									</div>
								) : null}
							</TabsContent>
						)
					})}
				</Tabs>
			</div>
		</ScrollArea>
	)
}
