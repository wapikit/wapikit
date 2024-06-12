'use client'

import BreadCrumb from '~/components/breadcrumb'
import { ContactListTableColumns } from '~/components/tables/columns'
import { TableComponent } from '~/components/tables/table'
import { buttonVariants } from '~/components/ui/button'
import { Heading } from '~/components/ui/heading'
import { Separator } from '~/components/ui/separator'
import { useGetContactLists, type ContactListSchema } from 'root/.generated'
import { Plus } from 'lucide-react'
import Link from 'next/link'
import { clsx } from 'clsx'
import { useSearchParams } from 'next/navigation'

const breadcrumbItems = [{ title: 'lists', link: '/lists' }]

const ListsPage = () => {
	// ! TODO:
	// * 1. Create a page for lists
	// * 2. Create a form to add a contact
	// * 3. Import bulk contact button
	// * 4. Bulk select actions : Export, Delete, Create a new List
	// * 5 . Individual contact actions : Edit, Delete, Add to List

	const searchParams = useSearchParams()

	const page = Number(searchParams.get('page') || 1)
	const pageLimit = Number(searchParams.get('limit') || 0) || 10
	// const offset = (page - 1) * pageLimit

	const contactResponse = useGetContactLists({})

	const totalUsers = contactResponse.data?.paginationMeta?.total || 0
	const pageCount = Math.ceil(totalUsers / pageLimit)
	const lists: ContactListSchema[] = contactResponse.data?.lists || []

	return (
		<>
			<div className="flex-1 space-y-4  p-4 pt-6 md:p-8">
				<BreadCrumb items={breadcrumbItems} />

				<div className="flex items-start justify-between">
					<Heading title={`Lists (${totalUsers})`} description="Manage lists" />

					<Link
						href={'/lists/new'}
						className={clsx(buttonVariants({ variant: 'default' }))}
					>
						<Plus className="mr-2 h-4 w-4" /> Add New
					</Link>
				</div>
				<Separator />

				<TableComponent
					searchKey="country"
					pageNo={page}
					columns={ContactListTableColumns}
					totalUsers={totalUsers}
					data={lists}
					pageCount={pageCount}
				/>
			</div>
		</>
	)
}

export default ListsPage
