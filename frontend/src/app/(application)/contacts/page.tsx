'use client'

import BreadCrumb from '~/components/breadcrumb'
import { ContactTableColumns } from '~/components/tables/columns'
import { TableComponent } from '~/components/tables/table'
import { buttonVariants } from '~/components/ui/button'
import { Heading } from '~/components/ui/heading'
import { Separator } from '~/components/ui/separator'
import { useGetContacts, type ContactSchema } from 'root/.generated'
import { Plus } from 'lucide-react'
import Link from 'next/link'
import { clsx } from 'clsx'
import { useSearchParams } from 'next/navigation'

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

	return (
		<>
			<div className="flex-1 space-y-4  p-4 pt-6 md:p-8">
				<BreadCrumb items={breadcrumbItems} />

				<div className="flex items-start justify-between">
					<Heading title={`Contacts (${totalUsers})`} description="Manage contacts" />

					<Link
						href={'/contacts/new'}
						className={clsx(buttonVariants({ variant: 'default' }))}
					>
						<Plus className="mr-2 h-4 w-4" /> Add New
					</Link>
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
