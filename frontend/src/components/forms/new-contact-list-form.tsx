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
import { zodResolver } from '@hookform/resolvers/zod'
import { useRouter } from 'next/navigation'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { type z } from 'zod'
import { successNotification, errorNotification } from '~/reusable-functions'
import { useCreateList, useUpdateListById, type ContactListSchema } from 'root/.generated'
import { Textarea } from '../ui/textarea'
import { NewContactListFormSchema } from '~/schema'
import { Icons } from '../icons'
import { MultiSelect } from '../multi-select'
import { useLayoutStore } from '~/store/layout.store'

interface FormProps {
	initialData: ContactListSchema | null
}

const NewContactListForm: React.FC<FormProps> = ({ initialData }) => {
	const router = useRouter()
	const { writeProperty, tags } = useLayoutStore()
	const [loading, setLoading] = useState(false)
	const action = initialData ? 'Save changes' : 'Create'
	const createLists = useCreateList()
	const updateList = useUpdateListById()

	const defaultValues = initialData
		? {
				...initialData
			}
		: {
				name: '',
				description: '',
				tags: []
			}

	const form = useForm<z.infer<typeof NewContactListFormSchema>>({
		resolver: zodResolver(NewContactListFormSchema),
		defaultValues
	})

	// eslint-disable-next-line @typescript-eslint/no-unused-vars
	const onSubmit = async (data: z.infer<typeof NewContactListFormSchema>) => {
		try {
			setLoading(true)
			if (initialData) {
				const response = await updateList.mutateAsync({
					id: initialData.uniqueId,
					data: {
						name: data.name,
						tags: [],
						description: data.description
					}
				})

				if (response.list.uniqueId) {
					successNotification({
						message: 'List updated successfully.'
					})
				} else {
					errorNotification({
						message: 'There was a problem with your request.'
					})
				}
			} else {
				const response = await createLists.mutateAsync(
					{
						data: {
							name: data.name,
							tags: [],
							description: data.description
						}
					},
					{
						onError() {
							errorNotification({
								message: 'There was a problem with your request.'
							})
						}
					}
				)

				if (response.list.uniqueId) {
					successNotification({
						message: 'List created successfully.'
					})
				} else {
					errorNotification({
						message: 'There was a problem with your request.'
					})
				}
			}
			router.push(`/lists`)
		} catch (error: any) {
			errorNotification({
				message: 'There was a problem with your request.'
			})
		} finally {
			setLoading(false)
		}
	}

	return (
		<>
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
											placeholder="List name"
											autoComplete="off"
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
											placeholder="List description"
											{...field}
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name="tagIds"
							render={({}) => (
								<FormItem className="tablet:w-3/4 tablet:gap-2 desktop:w-1/2 flex flex-col gap-1 ">
									<FormLabel>Select the tags to add</FormLabel>
									<MultiSelect
										options={
											tags.map(tag => ({
												label: tag.label,
												value: tag.uniqueId
											})) || []
										}
										onValueChange={e => {
											form.setValue('tagIds', e, {
												shouldValidate: true
											})
										}}
										defaultValue={form.watch('tagIds')}
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
					</div>
					<div className="flex w-fit flex-row gap-3 ">
						<Button disabled={loading} className="ml-auto" type="submit">
							{action}
						</Button>
					</div>
				</form>
			</Form>
		</>
	)
}

export default NewContactListForm
