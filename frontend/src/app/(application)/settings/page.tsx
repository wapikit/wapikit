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
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '~/components/ui/tooltip'
import { errorNotification, materialConfirm, successNotification } from '~/reusable-functions'
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue
} from '~/components/ui/select'
import { Heading } from '~/components/ui/heading'
import { Separator } from '~/components/ui/separator'
import { Plus, EyeIcon, Clipboard } from 'lucide-react'
import { useCopyToClipboard } from 'usehooks-ts'
import {
	RolePermissionEnum,
	useCreateOrganizationRole,
	useGetOrganizationRoleById,
	regenerateApiKey as regenerateApiKeyQuery,
	getApiKeys as getApiKeysQuery,
	useUpdateOrganizationRoleById,
	useUpdateWhatsappBusinessAccountDetails,
	getAllPhoneNumbers,
	useGetOrganizationRoles,
	useUpdateUser,
	useUpdateOrganization,
	AiModelEnum,
	getAIConfiguration
} from 'root/.generated'
import { Modal } from '~/components/ui/modal'
import {
	EmailNotificationConfigurationFormSchema,
	NewRoleFormSchema,
	OrganizationAiModelConfigurationSchema,
	OrganizationUpdateFormSchema,
	SlackNotificationConfigurationFormSchema,
	UserUpdateFormSchema,
	WhatsappBusinessAccountDetailsFormSchema
} from '~/schema'
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
import { useAuthState } from '~/hooks/use-auth-state'
import LoadingSpinner from '~/components/loader'
import { Textarea } from '~/components/ui/textarea'
import DocumentationPitch from '~/components/forms/documentation-pitch'
import { Switch } from '~/components/ui/switch'
import { Icons } from '~/components/icons'
import dynamic from 'next/dynamic'

export default function SettingsPage() {
	const { user, isOwner, currentOrganization, writeProperty, phoneNumbers, featureFlags } =
		useLayoutStore()

	enum SettingTabEnum {
		Account = 'account',
		Organization = 'organization',
		WhatsAppBusinessAccount = 'whatsapp-business-account',
		ApiKey = 'api-key',
		Rbac = 'rbac',
		Notifications = 'notifications',
		AiSettings = 'ai-settings',
		SubscriptionSettings = 'subscription-settings'
	}

	const tabs = [
		{
			slug: SettingTabEnum.Account,
			title: 'Account'
		},
		{
			slug: SettingTabEnum.Organization,
			title: 'Organization'
		},
		{
			slug: SettingTabEnum.WhatsAppBusinessAccount,
			title: 'WhatsApp Settings'
		},
		...(featureFlags?.SystemFeatureFlags.isApiAccessEnabled
			? [
					{
						slug: SettingTabEnum.ApiKey,
						title: 'API Key'
					}
				]
			: []),
		...(featureFlags?.SystemFeatureFlags.isRoleBasedAccessControlEnabled
			? [
					{
						slug: SettingTabEnum.Rbac,
						title: 'Access Control (RBAC)'
					}
				]
			: []),
		...(featureFlags?.SystemFeatureFlags.isAiIntegrationEnabled
			? [
					{
						slug: SettingTabEnum.AiSettings,
						title: 'AI Settings'
					}
				]
			: []),
		...(!featureFlags?.SystemFeatureFlags.isCloudEdition
			? [
					{
						slug: SettingTabEnum.Notifications,
						title: 'Notifications'
					}
				]
			: []),
		...(featureFlags?.SystemFeatureFlags.isCloudEdition
			? [
					{
						slug: SettingTabEnum.SubscriptionSettings,
						title: 'Subscription Settings'
					}
				]
			: [])
	]

	const SubscriptionSettings = dynamic(
		() =>
			process.env.NEXT_PUBLIC_IS_MANAGED_CLOUD_EDITION === 'true'
				? import('~/enterprise/components/settings/subscription')
				: Promise.resolve(() => null),
		{ ssr: false }
	)

	const searchParams = useSearchParams()
	const router = useRouter()
	const rolesDataSetRef = useRef(false)
	const { authState } = useAuthState()

	const copyToClipboard = useCopyToClipboard()[1]

	const page = Number(searchParams.get('page') || 1)
	const pageLimit = Number(searchParams.get('limit') || 0) || 10

	const [apiKey, setApiKey] = useState<string | null>(null)

	const [whatsAppBusinessAccountDetailsVisibility, setWhatsAppBusinessAccountDetailsVisibility] =
		useState({
			whatsappBusinessAccountId: false,
			apiToken: false,
			webhookSecret: false
		})

	const [isRoleCreationModelOpen, setIsRoleCreationModelOpen] = useState(false)
	const [roleIdToEdit, setRoleIdToEdit] = useState<string | null>(null)
	const [activeTab, setActiveTab] = useState(
		searchParams.get('tab')?.toString() || SettingTabEnum.Account
	)
	const [isBusy, setIsBusy] = useState(false)

	const [showWebhookSecret, setShowWebhookSecret] = useState(false)
	const [showAiApiKey, setShowAiApiKey] = useState(false)

	const createRoleMutation = useCreateOrganizationRole()
	const updateRoleMutation = useUpdateOrganizationRoleById()
	const updateWhatsappBusinessAccountDetailsMutation = useUpdateWhatsappBusinessAccountDetails()
	const updateUserMutation = useUpdateUser()
	const updateOrganizationMutation = useUpdateOrganization()
	const { data: roleData } = useGetOrganizationRoleById('', {
		query: {
			enabled: !!roleIdToEdit
		}
	})
	const { data: rolesResponse, refetch: refetchRoles } = useGetOrganizationRoles({
		page: page || 1,
		per_page: pageLimit || 10
	})

	const newRoleForm = useForm<z.infer<typeof NewRoleFormSchema>>({
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

	const userUpdateForm = useForm<z.infer<typeof UserUpdateFormSchema>>({
		resolver: zodResolver(UserUpdateFormSchema),
		defaultValues: {
			name: user?.name || ''
		}
	})

	const organizationUpdateForm = useForm<z.infer<typeof OrganizationUpdateFormSchema>>({
		resolver: zodResolver(OrganizationUpdateFormSchema),
		defaultValues: {
			name: currentOrganization?.name || '',
			description: currentOrganization?.description || ''
		}
	})

	const aiConfigurationForm = useForm<z.infer<typeof OrganizationAiModelConfigurationSchema>>({
		resolver: zodResolver(OrganizationAiModelConfigurationSchema),
		defaultValues: {
			isEnabled: currentOrganization?.aiConfiguration?.isEnabled ? true : false,
			model: currentOrganization?.aiConfiguration?.model || undefined,
			apiKey: ''
		}
	})

	const whatsappBusinessAccountIdForm = useForm<
		z.infer<typeof WhatsappBusinessAccountDetailsFormSchema>
	>({
		resolver: zodResolver(WhatsappBusinessAccountDetailsFormSchema),
		defaultValues: {
			whatsappBusinessAccountId: currentOrganization?.businessAccountId || undefined,
			apiToken: currentOrganization?.whatsappBusinessAccountDetails?.accessToken || undefined
		}
	})

	const slackNotificationConfigurationForm = useForm<
		z.infer<typeof SlackNotificationConfigurationFormSchema>
	>({
		resolver: zodResolver(SlackNotificationConfigurationFormSchema),
		defaultValues: {
			slackChannel:
				currentOrganization?.slackNotificationConfiguration?.slackChannel || undefined,
			slackWebhookUrl:
				currentOrganization?.slackNotificationConfiguration?.slackWebhookUrl || undefined
		}
	})

	const emailNotificationConfigurationForm = useForm<
		z.infer<typeof EmailNotificationConfigurationFormSchema>
	>({
		resolver: zodResolver(EmailNotificationConfigurationFormSchema),
		defaultValues: {
			smtpHost: currentOrganization?.emailNotificationConfiguration?.smtpHost || undefined,
			smtpPort: currentOrganization?.emailNotificationConfiguration?.smtpPort || undefined,
			smtpUsername:
				currentOrganization?.emailNotificationConfiguration?.smtpUsername || undefined,
			smtpPassword:
				currentOrganization?.emailNotificationConfiguration?.smtpPassword || undefined
		}
	})

	useEffect(() => {
		if (
			whatsappBusinessAccountIdForm.formState.touchedFields.apiToken ||
			whatsappBusinessAccountIdForm.formState.touchedFields.whatsappBusinessAccountId
		) {
			return
		}

		if (currentOrganization?.whatsappBusinessAccountDetails) {
			whatsappBusinessAccountIdForm.setValue(
				'whatsappBusinessAccountId',
				currentOrganization.whatsappBusinessAccountDetails.businessAccountId,
				{
					shouldTouch: true
				}
			)

			whatsappBusinessAccountIdForm.setValue(
				'apiToken',
				currentOrganization.whatsappBusinessAccountDetails.accessToken,
				{
					shouldTouch: true
				}
			)
		}

		return () => {
			if (
				whatsappBusinessAccountIdForm.formState.isDirty &&
				!whatsappBusinessAccountIdForm.formState.isSubmitting
			) {
				whatsappBusinessAccountIdForm.reset()
			}
		}
	}, [currentOrganization?.whatsappBusinessAccountDetails, whatsappBusinessAccountIdForm])

	useEffect(() => {
		if (
			aiConfigurationForm.formState.touchedFields.model ||
			aiConfigurationForm.formState.touchedFields.isEnabled ||
			aiConfigurationForm.formState.touchedFields.apiKey
		) {
			return
		}

		if (currentOrganization?.aiConfiguration) {
			aiConfigurationForm.setValue('model', currentOrganization.aiConfiguration.model, {
				shouldTouch: true
			})

			aiConfigurationForm.setValue(
				'isEnabled',
				currentOrganization.aiConfiguration.isEnabled || false,
				{
					shouldTouch: true
				}
			)
		}

		return () => {
			if (
				aiConfigurationForm.formState.isDirty &&
				!aiConfigurationForm.formState.isSubmitting
			) {
				aiConfigurationForm.reset()
			}
		}
	}, [currentOrganization?.aiConfiguration, aiConfigurationForm])

	useEffect(() => {
		if (rolesDataSetRef.current) return
		if (roleData) {
			newRoleForm.reset({
				name: roleData.role.name,
				description: roleData.role.description,
				permissions: roleData.role.permissions
			})
			rolesDataSetRef.current = true
		}
	}, [roleData, newRoleForm])

	useEffect(() => {
		const tab = searchParams.get('tab') || 'account'
		if (tab) {
			setActiveTab(() => tab)
		}
	}, [searchParams])

	useEffect(() => {
		if (roleIdToEdit) {
			setIsRoleCreationModelOpen(true)
		}
	}, [roleIdToEdit])

	useEffect(() => {
		if (!userUpdateForm.formState.touchedFields.name) {
			userUpdateForm.setValue('name', user?.name || '', {
				shouldTouch: false,
				shouldDirty: false,
				shouldValidate: true
			})
		}
	}, [user?.name, userUpdateForm])

	async function deleteOrganization() {
		try {
			setIsBusy(true)
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
		} finally {
			setIsBusy(false)
		}
	}

	async function leaveOrganization() {
		try {
			setIsBusy(true)
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
		} finally {
			setIsBusy(false)
		}
	}

	async function updateOrganizationWhatsAppBusinessAccountDetails(
		data: z.infer<typeof WhatsappBusinessAccountDetailsFormSchema>
	) {
		try {
			if (!currentOrganization) return

			const response = await updateWhatsappBusinessAccountDetailsMutation.mutateAsync({
				data: {
					businessAccountId: data.whatsappBusinessAccountId,
					accessToken: data.apiToken
				}
			})

			if (response) {
				writeProperty({
					currentOrganization: {
						...currentOrganization,
						whatsappBusinessAccountDetails: {
							businessAccountId: response.businessAccountId,
							accessToken: response.accessToken,
							webhookSecret: response.webhookSecret
						}
					}
				})
			} else {
				errorNotification({
					message: 'Error updating WhatsApp Business Account ID'
				})
			}
		} catch (error) {
			console.error(error)
			errorNotification({
				message: 'Error updating WhatsApp Business Account ID'
			})
		}
	}

	async function updateOrganizationDetails(data: z.infer<typeof OrganizationUpdateFormSchema>) {
		try {
			if (!currentOrganization) {
				return
			}
			setIsBusy(true)
			const response = await updateOrganizationMutation.mutateAsync({
				data: {
					name: data.name,
					description: data.description
				},
				id: currentOrganization?.uniqueId
			})

			if (response.isUpdated) {
				writeProperty({
					currentOrganization: {
						...currentOrganization,
						name: data.name,
						description: data.description
					}
				})
				successNotification({
					message: 'Organization details updated successfully'
				})
			} else {
				errorNotification({
					message: 'Error updating organization details'
				})
			}
		} catch (error) {
			console.error(error)
			errorNotification({
				message: 'Error updating organization details'
			})
		} finally {
			setIsBusy(false)
		}
	}

	async function updateUserDetails(data: z.infer<typeof UserUpdateFormSchema>) {
		try {
			if (!user) return
			setIsBusy(true)

			const response = await updateUserMutation.mutateAsync({
				data: {
					name: data.name
				}
			})

			if (response.isUpdated) {
				writeProperty({
					user: {
						...user,
						name: data.name
					}
				})
				successNotification({
					message: 'User details updated successfully'
				})
			} else {
				errorNotification({
					message: 'Error updating user details'
				})
			}
		} catch (error) {
			console.error(error)
			errorNotification({
				message: 'Error updating user details'
			})
		} finally {
			setIsBusy(false)
		}
	}

	async function submitRoleForm(data: z.infer<typeof NewRoleFormSchema>) {
		try {
			setIsBusy(true)
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
					await refetchRoles()
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

	async function copyApiKey() {
		try {
			setIsBusy(true)
			const apiKey = await getApiKeysQuery()
			if (!apiKey) {
				errorNotification({
					message: 'Error copying API key'
				})
			} else {
				await copyToClipboard(apiKey.apiKey.key)
				successNotification({
					message: 'API key copied to clipboard'
				})
			}
		} catch (error) {
			console.error({
				message: 'Error copying API key'
			})

			errorNotification({
				message: 'Error copying API key'
			})
		} finally {
			setIsBusy(false)
		}
	}

	async function getApiKey() {
		try {
			setIsBusy(true)
			const apiKey = await getApiKeysQuery()
			if (!apiKey) {
				errorNotification({
					message: 'Error copying API key'
				})
			} else {
				setApiKey(apiKey.apiKey.key)
			}
		} catch (error) {
			console.error(error)
			errorNotification({
				message: 'Error getting API key'
			})
		} finally {
			setIsBusy(false)
		}
	}

	async function regenerateApiKey() {
		try {
			setIsBusy(true)
			const confirmation = await materialConfirm({
				title: 'Regenerate API Key',
				description:
					'Are you sure you want to regenerate the API key? This will invalidate the current key.'
			})

			if (!confirmation) {
				return
			}

			const response = await regenerateApiKeyQuery()

			if (response.apiKey) {
				successNotification({
					message: 'API key regenerated successfully'
				})

				await copyToClipboard(response.apiKey.key)
				setApiKey(response.apiKey.key)
				successNotification({
					message: 'API key copied to clipboard'
				})
			} else {
				errorNotification({
					message: 'Error regenerating API key'
				})
			}
		} catch (error) {
			console.error(error)
			errorNotification({
				message: 'Error regenerating API key'
			})
		} finally {
			setIsBusy(false)
		}
	}

	async function syncPhoneNumbers() {
		try {
			setIsBusy(true)
			const confirmed = await materialConfirm({
				title: 'Sync Phone Numbers',
				description:
					'Are you sure you want to sync phone numbers with WhatsApp Business Account? This need all current campaigns to be either completed or no paused campaigns.'
			})
			if (!confirmed) return

			const response = await getAllPhoneNumbers()

			if (response) {
				writeProperty({
					phoneNumbers: response
				})
				successNotification({
					message: 'Phone numbers synced successfully'
				})
			} else {
				errorNotification({
					message: 'Error syncing phone numbers'
				})
			}
		} catch (error) {
			console.error(error)
			errorNotification({
				message: 'Error syncing phone numbers'
			})
		} finally {
			setIsBusy(false)
		}
	}

	async function updateSlackNotificationConfiguration(
		data: z.infer<typeof SlackNotificationConfigurationFormSchema>
	) {
		try {
			if (isBusy || !currentOrganization) return

			setIsBusy(() => true)

			const response = await updateOrganizationMutation.mutateAsync({
				data: {
					...currentOrganization,
					aiConfiguration: undefined,
					slackNotificationConfiguration: {
						slackChannel: data.slackChannel,
						slackWebhookUrl: data.slackWebhookUrl
					}
				},
				id: currentOrganization.uniqueId
			})

			if (response.isUpdated) {
				writeProperty({
					currentOrganization: {
						...currentOrganization,
						slackNotificationConfiguration: {
							slackChannel: data.slackChannel,
							slackWebhookUrl: data.slackWebhookUrl
						}
					}
				})
				successNotification({
					message: 'Slack notification configuration updated successfully'
				})
			} else {
				errorNotification({
					message: 'Error updating slack notification configuration'
				})
			}
		} catch (error) {
			console.error(error)
			errorNotification({
				message: 'Error updating slack notification configuration'
			})
		} finally {
			setIsBusy(() => false)
		}
	}

	async function updateEmailNotificationConfiguration(
		data: z.infer<typeof EmailNotificationConfigurationFormSchema>
	) {
		try {
			if (isBusy || !currentOrganization) return

			setIsBusy(() => true)

			const response = await updateOrganizationMutation.mutateAsync({
				data: {
					...currentOrganization,
					aiConfiguration: undefined,
					emailNotificationConfiguration: {
						smtpHost: data.smtpHost,
						smtpPort: data.smtpPort,
						smtpUsername: data.smtpUsername,
						smtpPassword: data.smtpPassword
					}
				},
				id: currentOrganization.uniqueId
			})

			if (response.isUpdated) {
				writeProperty({
					currentOrganization: {
						...currentOrganization,
						emailNotificationConfiguration: {
							smtpHost: data.smtpHost,
							smtpPort: data.smtpPort,
							smtpUsername: data.smtpUsername,
							smtpPassword: data.smtpPassword
						}
					}
				})
				successNotification({
					message: 'Email notification configuration updated successfully'
				})
			} else {
				errorNotification({
					message: 'Error updating email notification configuration'
				})
			}
		} catch (error) {
			console.error(error)
			errorNotification({
				message: 'Error updating email notification configuration'
			})
		} finally {
			setIsBusy(() => false)
		}
	}

	async function updateAiConfiguration(
		data: z.infer<typeof OrganizationAiModelConfigurationSchema>
	) {
		try {
			if (isBusy || !currentOrganization) return

			setIsBusy(() => true)

			const response = await updateOrganizationMutation.mutateAsync({
				data: {
					...currentOrganization,
					aiConfiguration: {
						model: data.model,
						apiKey: data.apiKey,
						isEnabled: data.isEnabled
					}
				},
				id: currentOrganization.uniqueId
			})

			if (response.isUpdated) {
				writeProperty({
					currentOrganization: {
						...currentOrganization,
						aiConfiguration: {
							model: data.model
						}
					}
				})
				successNotification({
					message: 'AI configuration updated successfully'
				})
			} else {
				errorNotification({
					message: 'Error updating ai configuration'
				})
			}
		} catch (error) {
			errorNotification({
				message: (error as Error).message
			})
		} finally {
			setIsBusy(() => false)
		}
	}

	return (
		<ScrollArea className="h-full pr-8">
			<div className="flex-1 space-y-4 p-4 pt-6 md:p-8">
				<div className="flex items-center justify-between space-y-2">
					<h2 className="text-3xl font-bold tracking-tight">Settings</h2>
				</div>
				<Tabs value={activeTab} className="space-y-4">
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
								{tab.slug === SettingTabEnum.WhatsAppBusinessAccount ? (
									<div className="mr-auto flex max-w-4xl flex-col gap-5">
										<Form {...whatsappBusinessAccountIdForm}>
											<form
												onSubmit={whatsappBusinessAccountIdForm.handleSubmit(
													updateOrganizationWhatsAppBusinessAccountDetails
												)}
											>
												<Card className="flex flex-row">
													<div className="flex-1">
														<CardContent className="mt-4 flex w-full flex-col gap-3">
															<FormField
																control={
																	whatsappBusinessAccountIdForm.control
																}
																name="whatsappBusinessAccountId"
																render={({ field }) => (
																	<FormItem className="w-full">
																		<FormLabel>
																			WhatsApp Business
																			Account ID
																		</FormLabel>
																		<FormControl>
																			<div className="flex flex-row gap-2">
																				<Input
																					disabled={
																						isBusy
																					}
																					placeholder="whatsapp business account id"
																					{...field}
																					autoComplete="off"
																					type={
																						whatsAppBusinessAccountDetailsVisibility.whatsappBusinessAccountId
																							? 'text'
																							: 'password'
																					}
																				/>
																				<span
																					className="rounded-md border p-1 px-2"
																					onClick={() => {
																						setWhatsAppBusinessAccountDetailsVisibility(
																							data => ({
																								...data,
																								whatsappBusinessAccountId:
																									!data.whatsappBusinessAccountId
																							})
																						)
																					}}
																				>
																					<EyeIcon className="size-5" />
																				</span>
																			</div>
																		</FormControl>
																		<FormMessage />
																	</FormItem>
																)}
															/>
															<FormField
																control={
																	whatsappBusinessAccountIdForm.control
																}
																name="apiToken"
																render={({ field }) => (
																	<FormItem className="w-full">
																		<FormLabel>
																			API key
																		</FormLabel>
																		<FormControl>
																			<div className="flex flex-row gap-2">
																				<Input
																					disabled={
																						isBusy
																					}
																					placeholder="whatsapp business account api token"
																					{...field}
																					autoComplete="off"
																					type={
																						whatsAppBusinessAccountDetailsVisibility.apiToken
																							? 'text'
																							: 'password'
																					}
																				/>
																				<span
																					className="rounded-md border p-1 px-2"
																					onClick={() => {
																						setWhatsAppBusinessAccountDetailsVisibility(
																							data => ({
																								...data,
																								apiToken:
																									!data.apiToken
																							})
																						)
																					}}
																				>
																					<EyeIcon className="size-5" />
																				</span>
																			</div>
																		</FormControl>
																		<FormMessage />
																	</FormItem>
																)}
															/>

															<Button
																type="submit"
																className="ml-auto w-fit"
																disabled={
																	isBusy ||
																	!whatsappBusinessAccountIdForm
																		.formState.isDirty
																}
															>
																Update
															</Button>
														</CardContent>
													</div>
												</Card>
											</form>
										</Form>

										<Card className="min-w-4xl flex-1 border-none ">
											<CardHeader>
												<CardTitle>Your Unique Webhook Secret</CardTitle>
												<CardDescription>
													Use this webhook secret while verifying your
													webhook endpoint in the meta developer
													dashboard.
												</CardDescription>
											</CardHeader>
											<CardContent className="flex flex-row items-center gap-1">
												<Input
													className="w-fit truncate px-6 disabled:text-slate-600"
													value={
														showWebhookSecret
															? currentOrganization
																	?.whatsappBusinessAccountDetails
																	?.webhookSecret ||
																'***********************'
															: '***********************'
													}
													disabled
												/>
												<span>
													<Button
														onClick={() => {
															const secret =
																currentOrganization
																	?.whatsappBusinessAccountDetails
																	?.webhookSecret

															if (secret) {
																copyToClipboard(secret).catch(
																	error => console.error(error)
																)

																successNotification({
																	message:
																		'Secret copied to clipboard'
																})
															}
														}}
														className="ml-2 flex w-fit gap-1"
														variant={'secondary'}
														disabled={isBusy}
													>
														<Clipboard className="size-5" />
														Copy
													</Button>
												</span>
												<span>
													<Button
														onClick={() => {
															setShowWebhookSecret(data => !data)
														}}
														className="ml-2 flex w-fit gap-1"
														variant={'secondary'}
														disabled={isBusy}
													>
														<EyeIcon className="size-5" />
														Show
													</Button>
												</span>
											</CardContent>
										</Card>

										<div className="flex flex-row gap-5">
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
																	onClick={() => {
																		syncPhoneNumbers().catch(
																			error =>
																				console.error(error)
																		)
																	}}
																	disabled={!isOwner || isBusy}
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
														phoneNumbers?.[0]?.display_phone_number ||
														'no organizations'
													}
												>
													<SelectTrigger>
														<SelectValue placeholder="Select Phone Number" />
													</SelectTrigger>

													<SelectContent>
														{!phoneNumbers ||
														phoneNumbers.length === 0 ? (
															<SelectItem value={'empty'} disabled>
																No Phone Numbers.
															</SelectItem>
														) : (
															<>
																{phoneNumbers?.map(phoneNumber => (
																	<SelectItem
																		key={
																			phoneNumber.display_phone_number
																		}
																		value={
																			phoneNumber.display_phone_number
																		}
																	>
																		{
																			phoneNumber.display_phone_number
																		}
																	</SelectItem>
																))}
															</>
														)}
													</SelectContent>
												</Select>
											</CardContent>
										</Card>
									</div>
								) : tab.slug === SettingTabEnum.Account ? (
									<div className="mr-auto flex max-w-4xl flex-col gap-5">
										{authState.isAuthenticated ? (
											<>
												<Form {...userUpdateForm}>
													<form
														onSubmit={userUpdateForm.handleSubmit(
															updateUserDetails
														)}
													>
														<Card className="flex flex-col p-2">
															<div className="flex flex-row">
																<div className="flex-1">
																	<CardHeader>
																		<CardTitle>Name</CardTitle>
																	</CardHeader>
																	<CardContent>
																		<FormField
																			control={
																				userUpdateForm.control
																			}
																			name="name"
																			render={({ field }) => (
																				<FormItem>
																					<FormControl>
																						<Input
																							disabled={
																								isBusy
																							}
																							placeholder="name"
																							{...field}
																							autoComplete="off"
																							value={userUpdateForm.watch(
																								'name'
																							)}
																						/>
																					</FormControl>
																					<FormMessage />
																				</FormItem>
																			)}
																		/>
																	</CardContent>
																</div>
																<div className="tremor-Divider-root mx-auto my-6 flex items-center justify-between gap-3 text-tremor-default text-tremor-content dark:text-dark-tremor-content">
																	<div className="bg-tremor-border dark:bg-dark-tremor-border h-full w-[1px]"></div>
																</div>
																<div className="flex-1">
																	<CardHeader>
																		<CardTitle>Email</CardTitle>
																	</CardHeader>
																	<CardContent>
																		<FormField
																			name="email"
																			render={({ field }) => (
																				<FormItem>
																					<FormControl>
																						<Input
																							type="email"
																							disabled={
																								true
																							}
																							placeholder="email"
																							{...field}
																							autoComplete="off"
																							value={
																								authState
																									.data
																									.user
																									?.email ||
																								undefined
																							}
																						/>
																					</FormControl>
																					<FormMessage />
																				</FormItem>
																			)}
																		/>
																	</CardContent>
																</div>
															</div>
															<Button
																type="submit"
																disabled={
																	isBusy ||
																	!userUpdateForm.formState
																		.isDirty
																}
																className="ml-auto mr-6 w-fit"
															>
																Save
															</Button>
														</Card>
													</form>
												</Form>

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
																		disabled={isBusy}
																		className="flex gap-2"
																	>
																		<Icons.trash />
																		Delete Account
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
											</>
										) : (
											<LoadingSpinner />
										)}
									</div>
								) : tab.slug === SettingTabEnum.ApiKey ? (
									<div className="mr-auto flex flex-col gap-5">
										<Card className="min-w-4xl flex-1 border-none ">
											<CardHeader>
												<CardTitle>API Access Key</CardTitle>
												<CardDescription>
													Use this API key to authenticate wapikit API
													requests.
												</CardDescription>
											</CardHeader>
											<CardContent className="flex flex-row items-center gap-1">
												{/* ! TODO: show API key on hover of the input the full api key if present */}
												<Input
													className="w-fit truncate px-6 disabled:text-slate-600"
													value={apiKey || '***********************'}
													disabled
												/>
												<span>
													<Button
														onClick={() => {
															getApiKey().catch(error =>
																console.error(error)
															)
														}}
														className="ml-2 flex w-fit gap-1"
														variant={'secondary'}
														disabled={isBusy}
													>
														<EyeIcon className="size-5" />
														Show
													</Button>
												</span>

												<span>
													<Button
														onClick={() => {
															copyApiKey().catch(error =>
																console.error(error)
															)
														}}
														className="ml-2 flex w-fit gap-1"
														variant={'secondary'}
														disabled={isBusy}
													>
														<Clipboard className="size-5" />
														Copy
													</Button>
												</span>

												{/* regenerate button */}
												<Button
													onClick={() => {
														regenerateApiKey().catch(error =>
															console.error(error)
														)
													}}
													className="ml-auto w-fit"
													variant={'destructive'}
													disabled={isBusy}
												>
													Regenerate
												</Button>
											</CardContent>
										</Card>
										<DocumentationPitch type="api-key" />
									</div>
								) : tab.slug === SettingTabEnum.Rbac ? (
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
												newRoleForm.reset()
											}}
										>
											<div className="flex w-full items-center justify-end space-x-2 pt-6">
												<Form {...newRoleForm}>
													<form
														onSubmit={newRoleForm.handleSubmit(
															submitRoleForm
														)}
														className="w-full space-y-8"
													>
														<div className="flex flex-col gap-8">
															<FormField
																control={newRoleForm.control}
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
																control={newRoleForm.control}
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
																control={newRoleForm.control}
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
																				newRoleForm.setValue(
																					'permissions',
																					e as RolePermissionEnum[],
																					{
																						shouldValidate:
																							true
																					}
																				)
																			}}
																			defaultValue={newRoleForm.watch(
																				'permissions'
																			)}
																			placeholder="Select permissions"
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
															Create Role
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
													disabled={isBusy}
												>
													<Plus className="mr-2 h-4 w-4" /> Add New
												</Button>
											</div>
										</div>
										<Separator />
										<RolesTable
											setRoleToEditId={setRoleIdToEdit}
											rolesResponse={rolesResponse}
										/>
									</div>
								) : tab.slug === SettingTabEnum.Organization ? (
									<div className="mr-auto flex max-w-4xl flex-col gap-5">
										{/* organization name update button */}

										<Form {...organizationUpdateForm}>
											<form
												onSubmit={organizationUpdateForm.handleSubmit(
													updateOrganizationDetails
												)}
												className="w-full space-y-8"
											>
												<Card className="flex flex-col p-2">
													<div className="flex flex-col">
														<div>
															<CardHeader>
																<CardTitle>
																	Organization Name
																</CardTitle>
															</CardHeader>
															<CardContent>
																<FormField
																	control={
																		organizationUpdateForm.control
																	}
																	name="name"
																	render={({ field }) => (
																		<FormItem>
																			<FormControl>
																				<Input
																					placeholder="default organization"
																					{...field}
																				/>
																			</FormControl>
																			<FormMessage />
																		</FormItem>
																	)}
																/>
															</CardContent>
														</div>
														<div>
															<CardHeader>
																<CardTitle>
																	Organization Description
																</CardTitle>
															</CardHeader>
															<CardContent>
																<FormField
																	control={
																		organizationUpdateForm.control
																	}
																	name="description"
																	render={({ field }) => (
																		<FormItem>
																			<FormControl>
																				<Textarea
																					placeholder="description..."
																					{...field}
																				/>
																			</FormControl>
																			<FormMessage />
																		</FormItem>
																	)}
																/>
															</CardContent>
														</div>
													</div>

													<Button
														disabled={
															isBusy ||
															!organizationUpdateForm.formState
																.isDirty
														}
														className="ml-auto mr-6 w-fit "
													>
														Save
													</Button>
												</Card>
											</form>
										</Form>

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
																	disabled={isOwner || isBusy}
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
																	disabled={isBusy}
																	onClick={() => {
																		deleteOrganization().catch(
																			error =>
																				console.error(error)
																		)
																	}}
																	className="flex gap-2"
																>
																	<Icons.trash />
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
								) : tab.slug === SettingTabEnum.Notifications ? (
									<>
										<div className="mr-auto flex max-w-4xl flex-col gap-5">
											{/* slack configuration */}

											<Form {...slackNotificationConfigurationForm}>
												<form
													onSubmit={slackNotificationConfigurationForm.handleSubmit(
														updateSlackNotificationConfiguration
													)}
													className="w-full space-y-8"
												>
													<Card className="flex flex-col p-2">
														<div className="flex flex-col">
															<div>
																<CardHeader>
																	<CardTitle>
																		Slack Notification
																		Configuration
																	</CardTitle>
																</CardHeader>
																<CardContent className="flex flex-col gap-3">
																	<FormField
																		control={
																			slackNotificationConfigurationForm.control
																		}
																		name="slackWebhookUrl"
																		render={({ field }) => (
																			<FormItem>
																				<FormLabel>
																					Webhook URL
																				</FormLabel>
																				<FormControl>
																					<Input
																						placeholder="https://hooks.slack.com/services/..."
																						{...field}
																					/>
																				</FormControl>
																				<FormMessage />
																			</FormItem>
																		)}
																	/>
																	<FormField
																		control={
																			slackNotificationConfigurationForm.control
																		}
																		name="slackChannel"
																		render={({ field }) => (
																			<FormItem>
																				<FormLabel>
																					Webhook Channel
																				</FormLabel>
																				<FormControl>
																					<Input
																						placeholder="#wapikit-notifications"
																						{...field}
																					/>
																				</FormControl>
																				<FormMessage />
																			</FormItem>
																		)}
																	/>
																</CardContent>
															</div>
														</div>

														<Button
															disabled={
																isBusy ||
																!slackNotificationConfigurationForm
																	.formState.isDirty
															}
															className="ml-auto mr-6 w-fit "
														>
															Save
														</Button>
													</Card>
												</form>
											</Form>
										</div>

										{/* email configuration */}
										<div className="mr-auto flex max-w-4xl flex-col gap-5">
											<Form {...emailNotificationConfigurationForm}>
												<form
													onSubmit={emailNotificationConfigurationForm.handleSubmit(
														updateEmailNotificationConfiguration
													)}
													className="w-full space-y-8"
												>
													<Card className="flex flex-col p-2">
														<div className="flex flex-col">
															<div>
																<CardHeader>
																	<CardTitle>
																		Email Notification
																		Configuration (SMTP)
																	</CardTitle>
																</CardHeader>
																<CardContent className="flex flex-col gap-3">
																	<FormField
																		control={
																			emailNotificationConfigurationForm.control
																		}
																		name="smtpHost"
																		render={({ field }) => (
																			<FormItem>
																				<FormLabel>
																					host
																				</FormLabel>
																				<FormControl>
																					<Input
																						placeholder="smtp.yoursite.com"
																						{...field}
																					/>
																				</FormControl>
																				<FormMessage />
																			</FormItem>
																		)}
																	/>
																	<FormField
																		control={
																			emailNotificationConfigurationForm.control
																		}
																		name="smtpUsername"
																		render={({ field }) => (
																			<FormItem>
																				<FormLabel>
																					Username
																				</FormLabel>
																				<FormControl>
																					<Input
																						placeholder="SMTP Username"
																						{...field}
																					/>
																				</FormControl>
																				<FormMessage />
																			</FormItem>
																		)}
																	/>
																	<FormField
																		control={
																			emailNotificationConfigurationForm.control
																		}
																		name="smtpPassword"
																		render={({ field }) => (
																			<FormItem>
																				<FormLabel>
																					Password
																				</FormLabel>
																				<FormControl>
																					<Input
																						type="SMTP Password"
																						placeholder="********"
																						{...field}
																					/>
																				</FormControl>
																				<FormMessage />
																			</FormItem>
																		)}
																	/>
																	<FormField
																		control={
																			emailNotificationConfigurationForm.control
																		}
																		name="smtpPort"
																		render={({ field }) => (
																			<FormItem>
																				<FormLabel>
																					Port
																				</FormLabel>
																				<FormControl>
																					<Input
																						placeholder="587"
																						{...field}
																					/>
																				</FormControl>
																				<FormMessage />
																			</FormItem>
																		)}
																	/>
																</CardContent>
															</div>
														</div>

														<Button
															disabled={
																isBusy ||
																!organizationUpdateForm.formState
																	.isDirty
															}
															className="ml-auto mr-6 w-fit "
														>
															Save
														</Button>
													</Card>
												</form>
											</Form>
										</div>
									</>
								) : tab.slug === SettingTabEnum.AiSettings ? (
									<>
										<div className="mr-auto flex max-w-4xl flex-col gap-5">
											<Form {...aiConfigurationForm}>
												<form
													onSubmit={aiConfigurationForm.handleSubmit(
														updateAiConfiguration
													)}
													className="w-full space-y-8"
												>
													<Card className="flex flex-col p-2">
														<div className="flex flex-col">
															<div>
																<CardHeader>
																	<CardTitle>
																		AI Model Configuration
																	</CardTitle>
																</CardHeader>
																<CardContent className="flex flex-col gap-3">
																	{/* enable/disable toggle */}
																	<FormField
																		control={
																			aiConfigurationForm.control
																		}
																		name="isEnabled"
																		render={({ field }) => (
																			<FormItem className="flex flex-row items-center gap-5">
																				<FormLabel>
																					Enable AI
																				</FormLabel>
																				<FormControl>
																					<Switch
																						checked={
																							field.value
																						}
																						onClick={() => {
																							aiConfigurationForm.setValue(
																								'isEnabled',
																								!field.value
																							)
																						}}
																						className="flex items-center !p-0"
																					/>
																				</FormControl>
																				<FormMessage />
																			</FormItem>
																		)}
																	/>

																	{/* model selector */}
																	<FormField
																		disabled={
																			!aiConfigurationForm.watch(
																				'isEnabled'
																			)
																		}
																		control={
																			aiConfigurationForm.control
																		}
																		name="model"
																		render={({ field }) => (
																			<FormItem>
																				<FormLabel>
																					AI Model
																				</FormLabel>
																				<Select
																					disabled={
																						isBusy ||
																						!aiConfigurationForm.watch(
																							'isEnabled'
																						)
																					}
																					onValueChange={
																						field.onChange
																					}
																					value={
																						field.value
																					}
																					defaultValue={
																						field.value
																					}
																				>
																					<FormControl>
																						<SelectTrigger>
																							<SelectValue
																								defaultValue={
																									field.value
																								}
																								placeholder="Choose AI Model"
																							/>
																						</SelectTrigger>
																					</FormControl>
																					<SelectContent>
																						{listStringEnumMembers(
																							AiModelEnum
																						).map(
																							status => {
																								return (
																									<SelectItem
																										key={
																											status.name
																										}
																										value={
																											status.value
																										}
																									>
																										{
																											status.name
																										}
																									</SelectItem>
																								)
																							}
																						)}
																					</SelectContent>
																				</Select>
																				<FormMessage />
																			</FormItem>
																		)}
																	/>

																	{/* api key */}
																	<FormField
																		disabled={
																			!aiConfigurationForm.watch(
																				'isEnabled'
																			)
																		}
																		control={
																			aiConfigurationForm.control
																		}
																		name="apiKey"
																		render={({ field }) => (
																			<FormItem>
																				<FormLabel>
																					API Key
																				</FormLabel>
																				<FormControl>
																					<div className="flex flex-row items-center gap-2">
																						<Input
																							autoComplete="off"
																							placeholder="*********"
																							{...field}
																						/>
																						<span>
																							<Button
																								onClick={e => {
																									e.preventDefault()

																									if (
																										showAiApiKey
																									) {
																										aiConfigurationForm.setValue(
																											'apiKey',
																											'',
																											{
																												shouldDirty:
																													false
																											}
																										)
																									} else {
																										getAIConfiguration()
																											.then(
																												data => {
																													aiConfigurationForm.setValue(
																														'apiKey',
																														data
																															.aiConfiguration
																															.apiKey,
																														{
																															shouldDirty:
																																false
																														}
																													)
																												}
																											)
																											.catch(
																												error =>
																													console.error(
																														error
																													)
																											)
																									}

																									setShowAiApiKey(
																										data =>
																											!data
																									)
																								}}
																								className="ml-2 flex w-fit gap-1"
																								variant={
																									'secondary'
																								}
																								disabled={
																									isBusy
																								}
																							>
																								<EyeIcon className="size-5" />
																							</Button>
																						</span>
																					</div>
																				</FormControl>
																				<FormMessage />
																			</FormItem>
																		)}
																	/>
																</CardContent>
															</div>
														</div>

														<Button
															disabled={
																isBusy ||
																!aiConfigurationForm.formState
																	.isDirty
															}
															className="ml-auto mr-6 w-fit "
														>
															Save
														</Button>
													</Card>
												</form>
											</Form>
										</div>
									</>
								) : tab.slug === SettingTabEnum.SubscriptionSettings &&
								  featureFlags?.SystemFeatureFlags.isCloudEdition ? (
									<SubscriptionSettings />
								) : null}
							</TabsContent>
						)
					})}
				</Tabs>
			</div>
		</ScrollArea>
	)
}
