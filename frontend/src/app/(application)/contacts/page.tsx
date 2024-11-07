'use client'

import BreadCrumb from '~/components/breadcrumb'
import { ContactTableColumns } from '~/components/tables/columns'
import { TableComponent } from '~/components/tables/table'
import { Button, buttonVariants } from '~/components/ui/button'
import { Heading } from '~/components/ui/heading'
import { Separator } from '~/components/ui/separator'
import {
	useDeleteContactById,
	useGetContactLists,
	useGetContacts,
	type ContactSchema
} from 'root/.generated'
import { ArrowDown, Plus } from 'lucide-react'
import Link from 'next/link'
import { clsx } from 'clsx'
import { useSearchParams } from 'next/navigation'
import { useForm } from 'react-hook-form'
import { BulkImportContactsFormSchema } from '~/schema'
import { type z } from 'zod'
import { useState } from 'react'
import { errorNotification, materialConfirm, successNotification } from '~/reusable-functions'
import { Modal } from '~/components/ui/modal'
import { zodResolver } from '@hookform/resolvers/zod'
import { FileUploaderComponent } from '~/components/file-uploader'
import {
	Form,
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage
} from '~/components/ui/form'
import { Input } from '~/components/ui/input'
import { useRouter } from 'next/navigation'
import { MultiSelect } from '~/components/multi-select'
import { AUTH_TOKEN_LS, BACKEND_URL } from '~/constants'

const breadcrumbItems = [{ title: 'Contacts', link: '/contacts' }]

const ContactsPage = () => {
	const searchParams = useSearchParams()
	const router = useRouter()
	const deleteContactByIdMutation = useDeleteContactById()

	const page = Number(searchParams.get('page') || 1)
	const pageLimit = Number(searchParams.get('limit') || 0) || 10
	const listIds = searchParams.get('lists')
	const status = searchParams.get('status')

	const { data: contactResponse, refetch: refetchContacts } = useGetContacts({
		...(listIds ? { list_id: listIds } : {}),
		...(status ? { status: status } : {}),
		page: page || 1,
		per_page: pageLimit || 10
	})

	const listsResponse = useGetContactLists({
		order: 'asc',
		page: 1,
		per_page: 50
	})

	const totalUsers = contactResponse?.paginationMeta?.total || 0
	const pageCount = Math.ceil(totalUsers / pageLimit)
	const contacts: ContactSchema[] = contactResponse?.contacts || []

	const [isBulkImportModalOpen, setIsBulkImportModalOpen] = useState(false)
	const [file, setFile] = useState<File | null>(null)

	const bulkImportForm = useForm<z.infer<typeof BulkImportContactsFormSchema>>({
		resolver: zodResolver(BulkImportContactsFormSchema)
	})

	const [isBulkImporting, setIsBulkImporting] = useState(false)

	async function onBulkContactImportFormSubmit(
		data: z.infer<typeof BulkImportContactsFormSchema>
	) {
		try {
			setIsBulkImporting(true)

			if (!file) {
				errorNotification({
					message: 'Please upload a file'
				})
				return
			}

			const formData = new FormData()
			formData.append('file', file)
			formData.append('delimiter', data.delimiter)
			formData.append('listIds', JSON.stringify(data.listIds)) // If `listIds` is an array, stringify it

			const response = await fetch(`${BACKEND_URL}/contacts/bulkImport`, {
				body: formData,
				method: 'POST',
				headers: {
					Accept: 'application/json',
					'x-access-token': localStorage.getItem(AUTH_TOKEN_LS) || ''
				},
				cache: 'no-cache',
				credentials: 'include'
			})

			const result = await response.json() // Assuming response is JSON
			console.log({ result })

			if (result.message) {
				successNotification({
					message: result.message
				})
				await refetchContacts()
				setIsBulkImportModalOpen(false)
			} else {
				errorNotification({
					message: 'Failed to import contacts'
				})
			}
		} catch (error) {
			console.error(error)
			errorNotification({
				message: 'An error occurred'
			})
		} finally {
			setIsBulkImporting(false)
		}
	}

	async function deleteContact(contactId: string) {
		try {
			if (!contactId) return

			const confirmation = await materialConfirm({
				title: 'Delete Contact',
				description: 'Are you sure you want to delete this contact?'
			})

			if (!confirmation) return

			const { data } = await deleteContactByIdMutation.mutateAsync({
				id: contactId
			})

			if (data) {
				successNotification({
					message: 'Contact deleted successfully'
				})
			} else {
				errorNotification({
					message: 'Failed to delete contact'
				})
			}
		} catch (error) {
			console.error('Error deleting contact', error)
			errorNotification({
				message: 'Error deleting contact'
			})
		}
	}

	return (
		<>
			{/* bulk import contacts */}
			<Modal
				title="Import Contacts"
				description="Upload a CSV file with the following columns: name, phoneNumber, attributes"
				isOpen={isBulkImportModalOpen}
				onClose={() => {
					setIsBulkImportModalOpen(false)
				}}
			>
				<div className="flex w-full items-center justify-end space-x-2 pt-6">
					<Form {...bulkImportForm}>
						<form
							onSubmit={bulkImportForm.handleSubmit(onBulkContactImportFormSubmit)}
							className="w-full space-y-8"
						>
							<div className="flex flex-col gap-8">
								<FormField
									control={bulkImportForm.control}
									name="file"
									render={({ field }) => (
										<FormItem>
											<FormLabel>Upload CSV File</FormLabel>
											<FileUploaderComponent
												descriptionString="CSV File"
												{...field}
												onFileUpload={e => {
													const file = e.target.files?.[0]
													console.log({ fileData: file?.name })

													if (!file) return
													setFile(() => file)
												}}
											/>
										</FormItem>
									)}
								/>

								<FormField
									control={bulkImportForm.control}
									name="delimiter"
									render={({ field }) => (
										<FormItem>
											<FormLabel>Delimiter</FormLabel>
											<FormControl>
												<Input
													disabled={isBulkImporting}
													placeholder="Column delimiter (e.g. ,)"
													{...field}
													autoComplete="off"
												/>
											</FormControl>
											<FormMessage />
										</FormItem>
									)}
								/>
								<FormField
									control={bulkImportForm.control}
									name="listIds"
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
													bulkImportForm.setValue('listIds', e, {
														shouldValidate: true
													})
												}}
												defaultValue={bulkImportForm.watch('listIds')}
												placeholder="Select lists"
												variant="default"
											/>
											<FormMessage />
										</FormItem>
									)}
								/>
							</div>
							<Button
								disabled={isBulkImporting}
								className="ml-auto mr-0 w-full"
								type="submit"
							>
								Import
							</Button>
						</form>
					</Form>
				</div>
			</Modal>

			<div className="flex-1 space-y-4  p-4 pt-6 md:p-8">
				<BreadCrumb items={breadcrumbItems} />

				<div className="flex items-start justify-between">
					<Heading title={`Contacts (${totalUsers})`} description="Manage contacts" />
					<div className="flex gap-2">
						<Button
							className={clsx(buttonVariants({ variant: 'default' }))}
							onClick={() => {
								setIsBulkImportModalOpen(true)
							}}
						>
							<ArrowDown className="mr-2 h-4 w-4" /> Import
						</Button>
						<Link
							href={'/contacts/new-or-edit'}
							className={clsx(buttonVariants({ variant: 'default' }))}
						>
							<Plus className="mr-2 h-4 w-4" /> Add New
						</Link>
					</div>
				</div>
				<Separator />

				<TableComponent
					searchKey="phone"
					pageNo={page}
					columns={ContactTableColumns}
					totalUsers={totalUsers}
					data={contacts}
					pageCount={pageCount}
					actions={[
						{
							icon: 'edit',
							label: 'Edit',
							onClick: (contactId: string) => {
								// redirect to the edit page with id in search param
								router.push(`/contacts/new-or-edit?id=${contactId}`)
							}
						},
						{
							icon: 'trash',
							label: 'Delete',
							onClick: (contactId: string) => {
								deleteContact(contactId).catch(console.error)
							}
						}
					]}
				/>
			</div>
		</>
	)
}

export default ContactsPage
