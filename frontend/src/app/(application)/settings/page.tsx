'use client'

import { Button } from '~/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '~/components/ui/card'
import { ScrollArea } from '~/components/ui/scroll-area'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '~/components/ui/tabs'
import { Input } from '~/components/ui/input'
import RolesTable from '~/components/settings/roles-table'
import { useEffect, useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { useLayoutStore } from '~/store/layout.store'
import { useSettingsStore } from '~/store/settings.store'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '~/components/ui/tooltip'
import { errorNotification, materialConfirm } from '~/reusable-functions'
import Image from 'next/image'
import { FileUploaderComponent } from '~/components/file-uploader'
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue
} from '~/components/ui/select'

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
			title: 'Access Control'
		}
	]

	const searchParams = useSearchParams()
	const router = useRouter()

	const [activeTab, setActiveTab] = useState(
		searchParams.get('tab')?.toString() || 'app-settings'
	)

	useEffect(() => {
		if (searchParams.get('tab')) {
			setActiveTab(searchParams.get('tab')?.toString() || 'account')
		}
	}, [searchParams])

	// const updateOrganizationSettings = useUpdateSettings()

	// async function handleSettingsUpdate() {
	// 	await updateOrganizationSettings.mutateAsync({
	// 		data: {}
	// 	})
	// }

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

	const { isOwner } = useLayoutStore()
	const { organizationSettings, whatsappSettings } = useSettingsStore()

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
									<RolesTable />
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
