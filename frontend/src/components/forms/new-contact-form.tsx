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
import { useRouter } from 'next/navigation'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { type z } from 'zod'
import { errorNotification, successNotification } from '~/reusable-functions'
import {
	useGetContactLists,
	useCreateContacts,
	useUpdateContactById,
	ContactStatusEnum,
	type ContactSchema
} from 'root/.generated'
import { Textarea } from '../ui/textarea'
import { NewContactFormSchema } from '~/schema'
import { listStringEnumMembers } from 'ts-enum-utils'
import { MultiSelect } from '../multi-select'

interface FormProps {
	initialData: ContactSchema | null
}

const NewContactForm: React.FC<FormProps> = ({ initialData }) => {
	const router = useRouter()
	const [loading, setLoading] = useState(false)

	const toastMessage = initialData ? 'Product updated.' : 'Product created.'
	const action = initialData ? 'Save changes' : 'Create'

	const listsResponse = useGetContactLists({
		order: 'asc',
		page: 1,
		per_page: 50
	})

	const createContact = useCreateContacts()
	const updateContact = useUpdateContactById()

	const defaultValues = initialData
		? {
				...initialData,
				lists: initialData.lists.map(list => list.uniqueId)
			}
		: {
				name: '',
				attributes: {},
				phone: '',
				lists: [],
				status: ContactStatusEnum.Active
			}

	const form = useForm<z.infer<typeof NewContactFormSchema>>({
		resolver: zodResolver(NewContactFormSchema),
		defaultValues
	})

	const onSubmit = async (data: z.infer<typeof NewContactFormSchema>) => {
		try {
			setLoading(true)
			if (initialData) {
				const response = await updateContact.mutateAsync({
					id: initialData.uniqueId,
					data: {
						name: data.name,
						attributes: data.attributes,
						phone: data.phone,
						status: data.status,
						lists: data.lists
					}
				})

				if (response.contact.uniqueId) {
					successNotification({
						message: toastMessage
					})
				} else {
					errorNotification({
						message: 'There was a problem with your request.'
					})
				}
			} else {
				const response = await createContact.mutateAsync(
					{
						data: [
							{
								name: data.name,
								attributes: data.attributes,
								phone: data.phone,
								status: data.status,
								listsIds: data.lists
							}
						]
					},
					{
						onError() {
							errorNotification({
								message: 'Something went wrong'
							})
						}
					}
				)

				if (response.message) {
					successNotification({
						message: response.message
					})
				} else {
					errorNotification({
						message: 'There was a problem with your request.'
					})
				}
			}
			router.push(`/contacts`)
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
											placeholder="Contact title"
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
											placeholder="Contact description"
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
							name="phone"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Phone Number</FormLabel>
									<FormControl>
										<Input
											disabled={loading}
											placeholder="Contact Phone Number"
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
							name="lists"
							render={({ field }) => (
								<FormItem className="tablet:w-3/4 tablet:gap-2 desktop:w-1/2 flex flex-col gap-1 ">
									<FormLabel>
										Select the Contact List to add this contact in.
									</FormLabel>
									{/* @ts-ignore */}
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
										// defaultValue={form.watch('lists')}
										placeholder="Select lists"
										variant="default"
										{...field}
									/>
									<FormMessage />
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="status"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Status</FormLabel>
									<Select
										disabled={loading}
										onValueChange={field.onChange}
										value={field.value}
										defaultValue={field.value}
									>
										<FormControl>
											<SelectTrigger>
												<SelectValue
													defaultValue={field.value}
													placeholder="Status"
												/>
											</SelectTrigger>
										</FormControl>
										<SelectContent>
											{listStringEnumMembers(ContactStatusEnum).map(
												status => {
													return (
														<SelectItem
															key={status.name}
															value={status.value}
														>
															{status.name}
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
					</div>
					<Button disabled={loading} className="ml-auto" type="submit">
						{action}
					</Button>
				</form>
			</Form>
		</>
	)
}

export default NewContactForm
