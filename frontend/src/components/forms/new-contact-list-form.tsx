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
// import {
// 	Select,
// 	SelectContent,
// 	SelectItem,
// 	SelectTrigger,
// 	SelectValue
// } from '~/components/ui/select'
import { zodResolver } from '@hookform/resolvers/zod'
import { Trash } from 'lucide-react'
import { useRouter } from 'next/navigation'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { type z } from 'zod'
import { successNotification, errorNotification } from '~/reusable-functions'
import {
	// useGetOrganizationTags,
	useCreateList,
	useUpdateListById,
	type ContactListSchema
} from 'root/.generated'
import { Textarea } from '../ui/textarea'
import { NewContactListFormSchema } from '~/schema'

interface FormProps {
	initialData: ContactListSchema | null
}

const NewContactListForm: React.FC<FormProps> = ({ initialData }) => {
	const router = useRouter()
	const [loading, setLoading] = useState(false)
	const action = initialData ? 'Save changes' : 'Create'

	// const tagsResponse = useGetOrganizationTags({
	// 	page: 1,
	// 	per_page: 50,
	// 	sortBy: 'asc'
	// })

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
							// tags: data.tags,
							tags: [],
							description: data.description
						}
					},
					{
						onError(error) {
							errorNotification({
								message: error.message || 'There was a problem with your request.'
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
											placeholder="List name"
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
						{/* <FormField
							control={form.control}
							name="tags"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Tags</FormLabel>
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
													placeholder="Add tags"
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
				</form>
			</Form>
		</>
	)
}

export default NewContactListForm
