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
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue
} from '~/components/ui/select'
import { zodResolver } from '@hookform/resolvers/zod'
import { Trash } from 'lucide-react'
import { useRouter } from 'next/navigation'
import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { type z } from 'zod'
import { errorNotification, materialConfirm, successNotification } from '~/reusable-functions'
import { NewCampaignSchema } from '~/schema'
import {
	type CampaignSchema,
	useCreateCampaign,
	useGetContactLists,
	useUpdateCampaignById,
	CampaignStatusEnum,
	useGetOrganizationTags,
	useGetAllPhoneNumbers,
	useGetAllTemplates,
	useDeleteCampaignById
} from 'root/.generated'
import { Textarea } from '../ui/textarea'
import { Checkbox } from '../ui/checkbox'
import { type CheckedState } from '@radix-ui/react-checkbox'
import { DatePicker } from '../ui/date-picker'
import { MultiSelect } from '../multi-select'
import { useLayoutStore } from '~/store/layout.store'
import { ReloadIcon } from '@radix-ui/react-icons'

interface FormProps {
	initialData: CampaignSchema | null
}

const NewCampaignForm: React.FC<FormProps> = ({ initialData }) => {
	const router = useRouter()
	const [loading, setLoading] = useState(false)
	const toastMessage = initialData ? 'Product updated.' : 'Product created.'
	const action = initialData ? 'Save changes' : 'Create'

	const [hasUnsavedChanges, setHasUnsavedChanges] = useState(false)
	const [isBusy, setIsBusy] = useState(false)
	const [isScheduled, setIsScheduled] = useState(initialData?.scheduledAt ? true : false)

	const { writeProperty } = useLayoutStore()

	const listsResponse = useGetContactLists({
		order: 'asc',
		page: 1,
		per_page: 50
	})

	const { data: phoneNumbersResponse, refetch: refetchPhoneNumbers } = useGetAllPhoneNumbers()
	const { data: templatesResponse, refetch: refetchMessageTemplates } = useGetAllTemplates()
	const { data: tags } = useGetOrganizationTags({
		page: 1,
		per_page: 50,
		sortBy: 'asc'
	})

	const createNewCampaign = useCreateCampaign()
	const deleteCampaignById = useDeleteCampaignById()
	const updateCampaign = useUpdateCampaignById()

	const defaultValues = initialData
		? {
				...initialData,
				tags: initialData.tags.map(tag => tag.uniqueId),
				lists: initialData.lists.map(list => list.uniqueId)
			}
		: {
				name: '',
				description: '',
				isLinkTrackingEnabled: false,
				lists: [],
				status: CampaignStatusEnum.Draft,
				tags: []
			}

	const form = useForm<CampaignFormValues>({
		resolver: zodResolver(NewCampaignSchema),
		defaultValues
	})

	const onSubmit = async (data: CampaignFormValues) => {
		try {
			console.log('submitting form......')
			console.log('data is', JSON.stringify(data, null, 4))
			setLoading(true)
			if (initialData) {
				console.log('updating campaign')
				const response = await updateCampaign.mutateAsync({
					id: initialData.uniqueId,
					data: {
						description: data.description,
						enableLinkTracking: data.isLinkTrackingEnabled,
						listIds: data.lists,
						name: data.name,
						// templateMessageId: data.templateId,
						tags: data.tags
					}
				})

				if (response.campaign) {
					successNotification({
						message: toastMessage
					})
				} else {
					errorNotification({
						message: 'There was a problem with your request.'
					})
				}
			} else {
				console.log('creating new campaign')
				const response = await createNewCampaign.mutateAsync({
					data: {
						description: data.description,
						isLinkTrackingEnabled: data.isLinkTrackingEnabled,
						listIds: data.lists,
						name: data.name,
						// templateMessageId: data.templateId,
						tags: data.tags
					}
				})

				if (response.campaign) {
					successNotification({
						message: toastMessage
					})
					router.push(`/dashboard/campaigns`)
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
			if (form.formState.isDirty) {
				setHasUnsavedChanges(true)
			} else if (form.formState.isSubmitted) {
				setHasUnsavedChanges(false)
			}
		}
	}, [form.formState.isDirty, form.formState.isSubmitted])

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
			<div className="flex flex-1 items-center justify-between">
				{initialData && (
					<Button
						disabled={loading}
						variant="destructive"
						size="sm"
						onClick={() => {
							deleteCampaign().catch(error => console.error(error))
						}}
					>
						<Trash className="h-4 w-4" />
					</Button>
				)}
			</div>
			<Form {...form}>
				<form
					onSubmit={e => {
						e.preventDefault()
						onSubmit(form.getValues()).catch(error => console.error(error))
					}}
					className="w-full space-y-8"
				>
					<div className="w-full space-y-8">
						<div className="flex flex-col gap-8">
							<FormField
								control={form.control}
								name="name"
								render={({ field }) => (
									<FormItem>
										<FormLabel>Name</FormLabel>
										<FormControl>
											<Input
												disabled={loading}
												placeholder="Campaign title"
												{...field}
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
								control={form.control}
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
												console.log({ e })
												form.setValue('lists', e, {
													shouldValidate: true
												})
											}}
											defaultValue={form.watch('lists')}
											placeholder="Select lists"
											variant="default"
										/>
										<FormMessage />
									</FormItem>
								)}
							/>

							<FormField
								control={form.control}
								name="tags"
								render={({}) => (
									<FormItem className="tablet:w-3/4 tablet:gap-2 desktop:w-1/2 flex flex-col gap-1 ">
										<FormLabel>Select the tags to add</FormLabel>
										<MultiSelect
											options={
												tags?.tags?.map(tag => ({
													label: tag.name,
													value: tag.uniqueId
												})) || []
											}
											onValueChange={e => {
												console.log({ e })
												form.setValue('tags', e, {
													shouldValidate: true
												})
											}}
											defaultValue={form.watch('tags')}
											placeholder="Select Tags"
											variant="default"
										/>
										<FormMessage />
									</FormItem>
								)}
							/>

							<FormField
								control={form.control}
								name="templateId"
								render={({ field }) => (
									<FormItem>
										<FormLabel className="flex flex-row items-center gap-2">
											Message Template
											<Button
												size={'sm'}
												variant={'secondary'}
												onClick={() => {
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
										<Select
											disabled={loading}
											onValueChange={field.onChange}
											// defaultValue={field.value}
										>
											<FormControl>
												<SelectTrigger>
													<SelectValue placeholder="Select message template" />
												</SelectTrigger>
											</FormControl>
											<SelectContent>
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
														{templatesResponse?.map(template => (
															<SelectItem
																key={template.id}
																value={template.name}
															>
																{template.name}
															</SelectItem>
														))}
													</>
												)}
											</SelectContent>
										</Select>
										<FormMessage />
									</FormItem>
								)}
							/>

							<FormField
								control={form.control}
								name="templateId"
								render={({ field }) => (
									<FormItem>
										<FormLabel className="flex flex-row items-center gap-2">
											Phone Number
											<Button
												size={'sm'}
												variant={'secondary'}
												onClick={() => {
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
										<Select
											disabled={loading}
											onValueChange={field.onChange}
											value={field.value || undefined}
										>
											<FormControl>
												<SelectTrigger>
													<SelectValue placeholder="Select Phone Numbers" />
												</SelectTrigger>
											</FormControl>
											<SelectContent>
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
																value={phone.display_phone_number}
															>
																{phone.display_phone_number}
															</SelectItem>
														))}
													</>
												)}
											</SelectContent>
										</Select>
										<FormMessage />
									</FormItem>
								)}
							/>

							<div className="flex items-center gap-6">
								<FormField
									control={form.control}
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
									control={form.control}
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
							{/* <FormField
							control={form.control}
							name="templateId"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Select Template</FormLabel>
									<Select
										disabled={loading}
										onValueChange={field.onChange}
										value={field.value}
										// defaultValue={field.value}
									>
										<FormControl>
											<SelectTrigger>
												<SelectValue
													defaultValue={field.value}
													placeholder="Select Template"
												/>
											</SelectTrigger>
										</FormControl>
										<SelectContent>
											{tagsResponse.data?.tags.map(tag => (
												<SelectItem key={tag.uniqueId} value={tag.uniqueId}>
													{tag.name}
												</SelectItem>
											))}
										</SelectContent>
									</Select>
									<FormMessage />
								</FormItem>
							)}
						/> */}
						</div>
						<Button disabled={loading} className="ml-auto" type="submit">
							{action}
						</Button>
					</div>
				</form>
			</Form>
		</>
	)
}

export default NewCampaignForm

type CampaignFormValues = z.infer<typeof NewCampaignSchema>
