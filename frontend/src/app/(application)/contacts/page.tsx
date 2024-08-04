'use client'

import BreadCrumb from '~/components/breadcrumb'
import { ContactTableColumns } from '~/components/tables/columns'
import { TableComponent } from '~/components/tables/table'
import { Button, buttonVariants } from '~/components/ui/button'
import { Heading } from '~/components/ui/heading'
import { Separator } from '~/components/ui/separator'
import { useGetContacts, type ContactSchema } from 'root/.generated'
import { ArrowDown, Plus } from 'lucide-react'
import Link from 'next/link'
import { clsx } from 'clsx'
import { useSearchParams } from 'next/navigation'
import { useForm } from 'react-hook-form'
import { BulkImportContactsFormSchema } from '~/schema'
import { type z } from 'zod'
import { useState } from 'react'
import { errorNotification } from '~/reusable-functions'
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

const breadcrumbItems = [{ title: 'Contacts', link: '/contacts' }]

const ContactsPage = () => {
	// ! TODO:
	// * 1. Create a page for contacts
	// * 2. Create a form to add a contact
	// * 3. Import bulk contact button
	// * 4. Bulk select actions : Export, Delete, Create a new List
	// * 5 . Individual contact actions : Edit, Delete, Add to List

	const searchParams = useSearchParams()

	const page = Number(searchParams.get('page') || 1)
	const pageLimit = Number(searchParams.get('limit') || 0) || 10
	const listIds = searchParams.get('lists')
	const status = searchParams.get('status')
	// const offset = (page - 1) * pageLimit

	const contactResponse = useGetContacts({
		...(listIds ? { list_id: listIds } : {}),
		...(status ? { status: status } : {}),
		page: page || 1,
		per_page: pageLimit || 10
	})

	const totalUsers = contactResponse.data?.paginationMeta?.total || 0
	const pageCount = Math.ceil(totalUsers / pageLimit)
	const contacts: ContactSchema[] = contactResponse.data?.contacts || []

	const [isBulkImportModalOpen, setIsBulkImportModalOpen] = useState(false)

	const defaultValues = {
		delimiter: '',
		file: null
	}

	const form = useForm<z.infer<typeof BulkImportContactsFormSchema>>({
		resolver: zodResolver(BulkImportContactsFormSchema),
		defaultValues
	})

	const [isBulkImporting, setIsBulkImporting] = useState(false)

	async function onSubmit() {
		try {
			setIsBulkImporting(true)

			// validate the
		} catch (error) {
			console.error(error)
			errorNotification({
				message: 'An error occurred'
			})
		} finally {
			setIsBulkImporting(false)
		}
	}

	return (
		<>
			{/* bulk import contacts */}
			<Modal
				title="Import Contacts"
				description="Upload a CSV file with the following columns: name, phoneNumber, email, attributes"
				isOpen={isBulkImportModalOpen}
				onClose={() => {
					setIsBulkImportModalOpen(false)
				}}
			>
				<div className="flex w-full items-center justify-end space-x-2 pt-6">
					<Form {...form}>
						<form onSubmit={form.handleSubmit(onSubmit)} className="w-full space-y-8">
							<div className="flex flex-col gap-8">
								<FileUploaderComponent
									descriptionString="CSV or a Zip File"
									onFileUpload={async e => {
										console.log({ e })
									}}
								/>
								<FormField
									control={form.control}
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
					searchKey="phoneNumber"
					pageNo={page}
					columns={ContactTableColumns}
					totalUsers={totalUsers}
					data={contacts}
					pageCount={pageCount}
				/>
			</div>
		</>
	)
}

export default ContactsPage
