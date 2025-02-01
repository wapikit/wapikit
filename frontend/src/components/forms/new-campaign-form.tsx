'use client'
import { Button } from '~/components/ui/button'
import {
	Form,
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage
} from '~/components/ui/form'
import { Input } from '~/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger } from '~/components/ui/select'
import { zodResolver } from '@hookform/resolvers/zod'
import { Pencil, Trash } from 'lucide-react'
import { useRouter } from 'next/navigation'
import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { type z } from 'zod'
import { errorNotification, materialConfirm, successNotification } from '~/reusable-functions'
import { NewCampaignSchema, TemplateComponentSchema } from '~/schema'
import {
	type CampaignSchema,
	useCreateCampaign,
	useGetContactLists,
	useUpdateCampaignById,
	useGetOrganizationTags,
	useGetAllPhoneNumbers,
	useGetAllTemplates,
	useDeleteCampaignById,
	getTemplateById,
	type UpdateCampaignSchema,
	CampaignStatusEnum
} from 'root/.generated'
import { Textarea } from '../ui/textarea'
import { Checkbox } from '../ui/checkbox'
import { type CheckedState } from '@radix-ui/react-checkbox'
import { DatePicker } from '../ui/date-picker'
import { MultiSelect } from '../multi-select'
import { useLayoutStore } from '~/store/layout.store'
import { ReloadIcon } from '@radix-ui/react-icons'
import { useAuthState } from '~/hooks/use-auth-state'
import * as React from 'react'
import {
	Drawer,
	DrawerContent,
	DrawerDescription,
	DrawerHeader,
	DrawerTitle
} from '~/components/ui/drawer'
import { Separator } from '../ui/separator'
import { ScrollArea } from '../ui/scroll-area'
import { isPresent } from 'ts-is-present'

import TemplateParameterForm from './template-parameter-form'
import TemplateMessageRenderer from '../chat/template-message-renderer'
import { Icons } from '../icons'

interface FormProps {
	initialData: CampaignSchema | null
}

const NewCampaignForm: React.FC<FormProps> = ({ initialData }) => {
	const router = useRouter()
	const { authState } = useAuthState()
	const [loading, setLoading] = useState(false)
	const toastMessage = initialData ? 'Campaign updated.' : 'Campaign created.'
	const action = initialData ? 'Save changes' : 'Create'

	const [newCampaignId, setNewCampaignId] = useState<string | null>(null)

	const [hasUnsavedChanges, setHasUnsavedChanges] = useState(false)
	const [isBusy, setIsBusy] = useState(false)
	const [isTemplateComponentsInputModalOpen, setIsTemplateComponentsInputModalOpen] =
		useState(false)
	const [isScheduled, setIsScheduled] = useState(initialData?.scheduledAt ? true : false)

	const { writeProperty } = useLayoutStore()

	const listsResponse = useGetContactLists({
		order: 'asc',
		page: 1,
		per_page: 50
	})

	const { data: phoneNumbersResponse, refetch: refetchPhoneNumbers } = useGetAllPhoneNumbers({
		query: {
			enabled: !!authState.isAuthenticated
		}
	})

	const { data: templatesResponse, refetch: refetchMessageTemplates } = useGetAllTemplates({
		query: {
			enabled: !!authState.isAuthenticated
		}
	})

	const { data: tags } = useGetOrganizationTags({
		page: 1,
		per_page: 50,
		sortBy: 'asc'
	})

	const createNewCampaign = useCreateCampaign()
	const deleteCampaignById = useDeleteCampaignById()
	const updateCampaign = useUpdateCampaignById()

	const campaignForm = useForm<CampaignFormValues>({
		resolver: zodResolver(NewCampaignSchema)
	})

	const templateMessageComponentParameterForm = useForm<z.infer<typeof TemplateComponentSchema>>({
		resolver: zodResolver(TemplateComponentSchema),
		defaultValues: {
			body: [],
			header: [],
			buttons: []
		}
	})

	useEffect(() => {
		if (initialData) {
			campaignForm.reset(
				{
					...initialData,
					name: initialData.name,
					description: initialData.description,
					tags: initialData.tags.map(tag => tag.uniqueId),
					lists: initialData.lists.map(list => list.uniqueId),
					templateId: initialData.templateMessageId,
					isLinkTrackingEnabled: initialData.isLinkTrackingEnabled,
					phoneNumberToUse: initialData.phoneNumberInUse
				},
				{
					keepDirty: false
				}
			)

			if (initialData.templateComponentParameters) {
				templateMessageComponentParameterForm.reset(
					{
						...initialData.templateComponentParameters
					},
					{
						keepDirty: true
					}
				)
			}
		}
	}, [campaignForm, initialData, templateMessageComponentParameterForm])

	const onSubmit = async (data: CampaignFormValues) => {
		try {
			setLoading(true)
			if (initialData) {
				const response = await updateCampaign.mutateAsync({
					id: initialData.uniqueId,
					data: {
						description: data.description,
						enableLinkTracking: data.isLinkTrackingEnabled,
						listIds: data.lists,
						name: data.name,
						templateMessageId: data.templateId,
						phoneNumber: data.phoneNumberToUse,
						tags: data.tags,
						status: initialData.status
					}
				})

				if (response.isUpdated) {
					successNotification({
						message: toastMessage
					})
				} else {
					errorNotification({
						message: 'There was a problem with your request.'
					})
				}
			} else {
				const response = await createNewCampaign.mutateAsync({
					data: {
						description: data.description,
						isLinkTrackingEnabled: data.isLinkTrackingEnabled,
						listIds: data.lists,
						name: data.name,
						templateMessageId: data.templateId,
						phoneNumberToUse: data.phoneNumberToUse,
						tags: data.tags
					}
				})

				if (response.campaign) {
					successNotification({
						message: toastMessage
					})
					setNewCampaignId(response.campaign.uniqueId)
					if (data.templateId) {
						// fetch the template here and show the modal

						const templateInuse = await getTemplateById(data.templateId)

						if (!templateInuse) {
							errorNotification({
								message:
									'Unable to fetch your selected message template. However, your campaign has been created successfully. You can edit it later.'
							})
						}

						setIsTemplateComponentsInputModalOpen(true)
					} else {
						router.push(`/campaigns`)
					}
				} else {
					errorNotification({
						message: 'There was a problem with your request.'
					})
				}
			}
		} catch (error: unknown) {
			errorNotification({
				message:
					error instanceof Error
						? error.message || 'There was a problem with your request.'
						: 'There was a problem with your request.'
			})
		} finally {
			setLoading(false)
		}
	}

	const handleTemplateComponentParameterSubmit = async (
		data: z.infer<typeof TemplateComponentSchema>
	) => {
		try {
			const campaignId = newCampaignId || initialData?.uniqueId

			if (!campaignId) {
				errorNotification({
					message: 'Something went wrong while creating the campaign.'
				})
				return
			}

			setLoading(true)
			const updateCampaignData: UpdateCampaignSchema = initialData
				? {
						...initialData,
						templateComponentParameters: data,
						enableLinkTracking: initialData.isLinkTrackingEnabled,
						listIds: initialData.lists.map(list => list.uniqueId),
						tags: initialData.tags.map(tag => tag.uniqueId),
						phoneNumber: initialData.phoneNumberInUse,
						templateMessageId: initialData.templateMessageId
					}
				: {
						description: campaignForm.getValues('description'),
						enableLinkTracking: campaignForm.getValues('isLinkTrackingEnabled'),
						listIds: campaignForm.getValues('lists'),
						name: campaignForm.getValues('name'),
						templateMessageId: campaignForm.getValues('templateId'),
						phoneNumber: campaignForm.getValues('phoneNumberToUse'),
						tags: campaignForm.getValues('tags'),
						status: CampaignStatusEnum.Draft
					}

			const response = await updateCampaign.mutateAsync({
				data: {
					...updateCampaignData
				},
				id: campaignId
			})

			if (response.isUpdated) {
				router.push(`/campaigns`)
			} else {
				errorNotification({
					message: 'Something went wrong while creating the campaign.'
				})
			}
		} catch (error) {
			console.error(error)
			errorNotification({
				message: 'Something went wrong while inviting the team member.'
			})
		} finally {
			setLoading(false)
		}
	}

	async function deleteCampaign() {
		try {
			setIsBusy(true)
			if (!initialData?.uniqueId) return

			const confirmation = await materialConfirm({
				title: 'Delete Campaign',
				description: 'Are you sure you want to delete this campaign?'
			})

			if (!confirmation) {
				return
			}

			const response = await deleteCampaignById.mutateAsync({
				id: initialData.uniqueId
			})

			if (response.data) {
				successNotification({
					message: 'Campaign deleted successfully.'
				})
				router.push(`/campaigns`)
			} else {
				errorNotification({
					message: 'Something went wrong while deleting the campaign.'
				})
			}
		} catch (error) {
			console.error(error)
			errorNotification({
				message: 'Something went wrong while deleting the campaign.'
			})
		} finally {
			setIsBusy(false)
		}
	}

	useEffect(() => {
		return () => {
			if (campaignForm.formState.isDirty) {
				setHasUnsavedChanges(true)
			} else if (campaignForm.formState.isSubmitted) {
				setHasUnsavedChanges(false)
			}
		}
	}, [campaignForm.formState.isDirty, campaignForm.formState.isSubmitted])

	useEffect(() => {
		function handleUnload(e: BeforeUnloadEvent) {
			if (hasUnsavedChanges) {
				e.preventDefault()
			}
		}

		// add a event listener to notify if the form has unsaved changes and user tries to leave the page
		window.addEventListener('beforeunload', handleUnload)

		return () => {
			window.removeEventListener('beforeunload', handleUnload)
		}
	}, [hasUnsavedChanges])

	return (
		<>
			<Drawer
				open={isTemplateComponentsInputModalOpen}
				dismissible={false}
				onClose={() => {
					// if in case template parameter has not been saved show a warning.
					const isDirty = templateMessageComponentParameterForm.formState.isDirty
					if (isDirty) {
						materialConfirm({
							title: 'Unsaved changes',
							description: 'You have unsaved changes. Are you sure you want to leave?'
						})
							.then((response: boolean) => {
								if (response) {
									setIsTemplateComponentsInputModalOpen(() => false)
									router.push(`/campaigns`)
								} else {
									// do not close, user  has clicked on cancel
									return
								}
							})
							.catch(error => {
								console.error(error)
							})
					} else {
						router.push(`/campaigns`)
					}
				}}
			>
				<DrawerContent className="max-h-[80vh] min-h-[80vh] px-10">
					<div className="mx-auto w-full">
						<DrawerHeader className="w-full">
							<DrawerTitle>Fill template components</DrawerTitle>
							<DrawerDescription>
								Add the values for the template components parameters. You may use
								templating variables to add dynamic values. For example, you can use
								first_name to add the first name of the contact. Check docs here.
							</DrawerDescription>
						</DrawerHeader>
						<Separator />
						<div className="flex w-full items-start justify-end space-x-2 pt-6">
							<div className="h-full flex-1">
								<ScrollArea className="h-full flex-1">
									<TemplateParameterForm
										handleTemplateComponentParameterSubmit={
											handleTemplateComponentParameterSubmit
										}
										isBusy={isBusy}
										setIsTemplateComponentsInputModalOpen={
											setIsTemplateComponentsInputModalOpen
										}
										templateMessageComponentParameterForm={
											templateMessageComponentParameterForm
										}
										key={'template-parameter-form'}
										template={templatesResponse?.find(template => {
											return (
												template.id === campaignForm.getValues('templateId')
											)
										})}
									/>
								</ScrollArea>
							</div>

							<Separator orientation="vertical" className="h-full" />

							<div className="h-full flex-1 rounded-md border">
								<div className="rounded-t-md bg-primary px-2 py-1 text-sm text-primary-foreground">
									Template Preview
								</div>
								<div className="relative h-full w-full rounded-b-md bg-[#ebe5de] p-4 dark:bg-[#202c33]">
									<div className='absolute inset-0 z-20 h-full w-full  bg-[url("/assets/conversations-canvas-bg.png")] bg-repeat opacity-20' />

									<div className="relative z-30 h-96">
										<TemplateMessageRenderer
											templateMessage={templatesResponse?.find(template => {
												return (
													template.id ===
													campaignForm.getValues('templateId')
												)
											})}
											parameterValues={templateMessageComponentParameterForm.getValues()}
										/>
									</div>
								</div>
							</div>
						</div>
					</div>
				</DrawerContent>
			</Drawer>

			<Form {...campaignForm}>
				<form
					onSubmit={e => {
						e.preventDefault()
						onSubmit(campaignForm.getValues()).catch(error => console.error(error))
					}}
					className="w-full space-y-8"
				>
					<div className="w-full space-y-8">
						<div className="flex flex-col gap-8">
							<FormField
								control={campaignForm.control}
								name="name"
								render={({ field }) => (
									<FormItem>
										<FormLabel>Name</FormLabel>
										<FormControl>
											<Input
												disabled={loading}
												placeholder="Campaign title"
												autoComplete="off"
												{...field}
											/>
										</FormControl>
										<FormMessage />
									</FormItem>
								)}
							/>
							<FormField
								control={campaignForm.control}
								name="description"
								render={({ field }) => (
									<FormItem>
										<FormLabel>Description</FormLabel>
										<FormControl>
											<Textarea
												disabled={loading}
												placeholder="Campaign description"
												{...field}
											/>
										</FormControl>
										<FormMessage />
									</FormItem>
								)}
							/>
							<FormField
								control={campaignForm.control}
								name="lists"
								render={({}) => (
									<FormItem className="tablet:w-3/4 tablet:gap-2 desktop:w-1/2 flex flex-col gap-1 ">
										<FormLabel>Select the lists</FormLabel>
										<MultiSelect
											options={
												listsResponse?.data?.lists.map(list => ({
													label: list.name,
													value: list.uniqueId
												})) || []
											}
											onValueChange={e => {
												campaignForm.setValue('lists', e, {
													shouldValidate: true
												})
											}}
											defaultValue={campaignForm.getValues('lists')}
											placeholder="Select lists"
											variant="default"
										/>
										<FormMessage />
									</FormItem>
								)}
							/>

							<FormField
								control={campaignForm.control}
								name="tags"
								render={({}) => (
									<FormItem className="tablet:w-3/4 tablet:gap-2 desktop:w-1/2 flex flex-col gap-1 ">
										<FormLabel>Select the tags to add</FormLabel>
										<MultiSelect
											options={
												tags?.tags?.map(tag => ({
													label: tag.label,
													value: tag.uniqueId
												})) || []
											}
											onValueChange={e => {
												campaignForm.setValue('tags', e, {
													shouldValidate: true
												})
											}}
											defaultValue={campaignForm.watch('tags')}
											placeholder="Select Tags"
											variant="default"
											showCloseButton={false}
											actionButtonConfig={{
												label: (
													<span className="flex items-center gap-2">
														<Icons.add className="h-4 w-4" />
														Create Tag
													</span>
												),
												onClick: () => {
													writeProperty({
														isCreateTagModalOpen: true
													})
												}
											}}
										/>
										<FormMessage />
									</FormItem>
								)}
							/>

							<FormField
								control={campaignForm.control}
								name="templateId"
								render={({ field }) => (
									<FormItem>
										<FormLabel className="flex flex-row items-center gap-2">
											Message Template
											<Button
												disabled={isBusy}
												size={'sm'}
												variant={'secondary'}
												type="button"
												onClick={e => {
													e.preventDefault()
													refetchMessageTemplates()
														.then(data => {
															writeProperty({
																templates: data.data || []
															})
														})
														.catch(error => console.error(error))
												}}
											>
												<ReloadIcon className="size-3" />
											</Button>
										</FormLabel>
										<FormControl>
											<Select
												disabled={loading}
												onValueChange={e => {
													field.onChange(e)
												}}
												name="templateId"
											>
												<SelectTrigger>
													<div>
														{templatesResponse
															?.map(template => {
																if (
																	template.id ===
																	campaignForm.getValues(
																		'templateId'
																	)
																) {
																	const stringToReturn = `${template.name} - ${template.language} - ${template.category}`
																	return stringToReturn
																} else {
																	return null
																}
															})
															.filter(isPresent)[0] ||
															'Select message template'}
													</div>
												</SelectTrigger>
												<SelectContent side="bottom" className="max-h-64">
													{!templatesResponse ||
													templatesResponse?.length === 0 ? (
														<SelectItem
															value={'no message template'}
															disabled
														>
															No message template.
														</SelectItem>
													) : (
														<>
															{templatesResponse?.map(
																(template, index) => (
																	<SelectItem
																		key={`${template.id}-${index}`}
																		value={template.id}
																	>
																		{template.name} -{' '}
																		{template.language} -{' '}
																		{template.category}
																	</SelectItem>
																)
															)}
														</>
													)}
												</SelectContent>
											</Select>
										</FormControl>
										<FormMessage />
									</FormItem>
								)}
							/>

							<FormField
								control={campaignForm.control}
								name="phoneNumberToUse"
								render={({ field }) => (
									<FormItem>
										<FormLabel className="flex flex-row items-center gap-2">
											Phone Number
											<Button
												disabled={isBusy}
												size={'sm'}
												type="button"
												variant={'secondary'}
												onClick={e => {
													e.preventDefault()
													refetchPhoneNumbers()
														.then(data => {
															writeProperty({
																phoneNumbers: data.data || []
															})
														})
														.catch(error => console.error(error))
												}}
											>
												<ReloadIcon className="size-3" />
											</Button>
										</FormLabel>
										<FormControl>
											<Select
												disabled={loading}
												onValueChange={field.onChange}
											>
												<SelectTrigger>
													<div>
														{phoneNumbersResponse
															?.map(phoneNumber => {
																if (
																	phoneNumber.id ===
																	campaignForm.getValues(
																		'phoneNumberToUse'
																	)
																) {
																	return phoneNumber.display_phone_number
																} else {
																	return null
																}
															})
															.filter(isPresent)[0] ||
															'Select Phone Number'}
													</div>
												</SelectTrigger>
												<SelectContent side="bottom" className="max-h-64">
													{!phoneNumbersResponse ||
													phoneNumbersResponse?.length === 0 ? (
														<SelectItem
															value={'no phone numbers found'}
															disabled
														>
															No Phone Numbers.
														</SelectItem>
													) : (
														<>
															{phoneNumbersResponse?.map(phone => (
																<SelectItem
																	key={phone.id}
																	value={phone.id}
																>
																	{phone.display_phone_number}
																</SelectItem>
															))}
														</>
													)}
												</SelectContent>
											</Select>
										</FormControl>
										<FormMessage />
									</FormItem>
								)}
							/>

							<div className="flex items-center gap-6">
								<FormField
									control={campaignForm.control}
									name="isLinkTrackingEnabled"
									render={({ field }) => (
										<FormItem className="flex items-center gap-2">
											<FormControl className="mt-2 flex items-center justify-center">
												<Checkbox
													disabled={loading}
													checked={field.value}
													onCheckedChange={field.onChange}
												/>
											</FormControl>
											<FormLabel>Enable Link Tracking</FormLabel>
											<FormMessage />
										</FormItem>
									)}
								/>
								<div className="">
									<FormItem className="flex items-center gap-2">
										<Checkbox
											className="mt-2"
											disabled={loading}
											checked={isScheduled}
											onCheckedChange={(e: CheckedState) => {
												setIsScheduled(() => !!e)
											}}
										/>
										<FormLabel>Schedule</FormLabel>
										<FormMessage />
									</FormItem>
								</div>
							</div>
							{isScheduled ? (
								<FormField
									control={campaignForm.control}
									name="schedule.date"
									render={({ field }) => (
										<FormItem>
											<DatePicker
												prefilledDate={
													field.value ? new Date(field.value) : undefined
												}
											/>
											<FormMessage />
										</FormItem>
									)}
								/>
							) : null}
						</div>

						<div className="sticky bottom-0 mr-auto flex w-full flex-1 items-start justify-start gap-2 bg-background py-5">
							<Button
								disabled={loading || isBusy || !campaignForm.formState.isDirty}
								className="ml-auto flex-1"
								type="submit"
							>
								{action}
							</Button>
							{initialData && (
								<>
									<Button
										disabled={
											loading ||
											isBusy ||
											!campaignForm.getValues('templateId')
										}
										variant="secondary"
										type="button"
										onClick={() => {
											setIsTemplateComponentsInputModalOpen(true)
										}}
										className="flex flex-1 items-center justify-center gap-1"
									>
										<Pencil className="h-4 w-4" />
										Edit Template Parameters
									</Button>

									<Button
										disabled={loading || isBusy}
										variant="destructive"
										type="button"
										onClick={() => {
											deleteCampaign().catch(error => console.error(error))
										}}
										className="flex flex-1 items-center justify-center gap-1"
									>
										<Trash className="h-4 w-4" />
										Delete
									</Button>
								</>
							)}
						</div>
					</div>
				</form>
			</Form>
		</>
	)
}

export default NewCampaignForm

type CampaignFormValues = z.infer<typeof NewCampaignSchema>
