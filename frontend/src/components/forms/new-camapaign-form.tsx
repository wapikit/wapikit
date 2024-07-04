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
import { errorNotification, successNotification } from '~/reusable-functions'
import { NewCampaignSchema } from '~/schema'
import {
	type CampaignSchema,
	useCreateCampaign,
	useGetContactLists,
	useUpdateCampaignById,
	CampaignStatusEnum,
	useGetOrganizationTags
} from 'root/.generated'
import { Textarea } from '../ui/textarea'
import { Checkbox } from '../ui/checkbox'
import { type CheckedState } from '@radix-ui/react-checkbox'
import { DatePicker } from '../ui/date-picker'

interface FormProps {
	initialData: CampaignSchema | null
}

const NewCampaignForm: React.FC<FormProps> = ({ initialData }) => {
	const router = useRouter()
	const [loading, setLoading] = useState(false)
	const toastMessage = initialData ? 'Product updated.' : 'Product created.'
	const action = initialData ? 'Save changes' : 'Create'

	const [hasUnsavedChanges, setHasUnsavedChanges] = useState(false)

	const [isScheduled, setIsScheduled] = useState(initialData?.scheduledAt ? true : false)

	const listsResponse = useGetContactLists({
		order: 'asc',
		page: 1,
		per_page: 50
	})

	const tagsResponse = useGetOrganizationTags({
		page: 1,
		per_page: 50,
		sortBy: 'asc'
	})

	const createNewCampaign = useCreateCampaign()
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
			setLoading(true)
			if (initialData) {
				const response = await updateCampaign.mutateAsync({
					id: initialData.uniqueId,
					data: {
						description: data.description,
						enableLinkTracking: true,
						listIds: data.lists,
						name: data.name,
						templateMessageId: data.templateId,
						tags: data.tags
					}
				})

				if (response.campaign.uniqueId) {
					toast({
						variant: 'default',
						title: 'Success!',
						description: toastMessage
					})
				} else {
					toast({
						variant: 'destructive',
						title: 'Uh oh! Something went wrong.',
						description: 'There was a problem with your request.'
					})
				}
			} else {
				const response = await createNewCampaign.mutateAsync(
					{
						data: {
							description: data.description,
							isLinkTrackingEnabled: true,
							listIds: data.lists,
							name: data.name,
							templateMessageId: data.templateId,
							tags: data.tags
						}
					},
					{
						onError(error) {
							toast({
								variant: 'destructive',
								title: 'Uh oh! Something went wrong.',
								description: error.message
							})
						}
					}
				)

				if (response.campaign.uniqueId) {
					successNotification({
						message: toastMessage
					})
				} else {
					errorNotification({
						message: 'There was a problem with your request.'
					})
				}
			}
			router.refresh()
			router.push(`/dashboard/products`)
			toast({
				variant: 'destructive',
				title: 'Uh oh! Something went wrong.',
				description: 'There was a problem with your request.'
			})
		} catch (error: any) {
			toast({
				variant: 'destructive',
				title: 'Uh oh! Something went wrong.',
				description: 'There was a problem with your request.'
			})
		} finally {
			setLoading(false)
		}
	}

	useEffect(() => {
		return () => {
			if (form.formState.isDirty) {
				setHasUnsavedChanges(true)
			}
		}
	}, [form.formState.isDirty])

	useEffect(() => {
		// add a event listener to notify if the form has unsaved changes and user tries to leave the page

		window.addEventListener('beforeunload', e => {
			if (hasUnsavedChanges) {
				e.preventDefault()
			}
		})
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
							// ! TODO: headless UI alert modal here
						}}
					>
						<Trash className="h-4 w-4" />
					</Button>
				)}
			</div>
			<Form {...form}>
				<form onSubmit={form.handleSubmit(onSubmit)} className="w-full space-y-8">
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
							render={({ field }) => (
								<FormItem>
									<FormLabel>Lists</FormLabel>
									<Select
										disabled={loading}
										onValueChange={field.onChange}
										value={field.value.join(',')}
										// defaultValue={field.value}
									>
										<FormControl>
											<SelectTrigger>
												<SelectValue
													defaultValue={field.value}
													placeholder="Select list"
												/>
											</SelectTrigger>
										</FormControl>
										<SelectContent>
											{!listsResponse.data?.lists ||
											listsResponse.data.lists.length === 0 ? (
												<SelectItem value={'no list'} disabled>
													No Lists created yet.
												</SelectItem>
											) : (
												<>
													{listsResponse.data?.lists.map(list => (
														<SelectItem
															key={list.uniqueId}
															value={list.uniqueId}
														>
															{list.name}
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
							name="tags"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Tags</FormLabel>
									<Select
										disabled={loading}
										onValueChange={field.onChange}
										value={field.value.join(',')}
										// defaultValue={field.value}
									>
										<FormControl>
											<SelectTrigger>
												<SelectValue
													defaultValue={field.value}
													placeholder="Add tags"
												/>
											</SelectTrigger>
										</FormControl>
										<SelectContent>
											{!tagsResponse.data?.tags ||
											tagsResponse.data.tags.length === 0 ? (
												<SelectItem value={'no list'} disabled>
													No Lists created yet.
												</SelectItem>
											) : (
												<>
													{tagsResponse.data?.tags.map(tag => (
														<SelectItem
															key={tag.uniqueId}
															value={tag.uniqueId}
														>
															{tag.name}
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
						<FormField
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
						/>
					</div>
					<Button disabled={loading} className="ml-auto" type="submit">
						{action}
					</Button>
				</form>
			</Form>
		</>
	)
}

export default NewCampaignForm

type CampaignFormValues = z.infer<typeof NewCampaignSchema>
