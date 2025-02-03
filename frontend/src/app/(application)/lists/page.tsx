'use client'

import BreadCrumb from '~/components/breadcrumb'
import { ContactListTableColumns } from '~/components/tables/columns'
import { TableComponent } from '~/components/tables/table'
import { buttonVariants } from '~/components/ui/button'
import { Heading } from '~/components/ui/heading'
import { Separator } from '~/components/ui/separator'
import { useDeleteListById, useGetContactLists, type ContactListSchema } from 'root/.generated'
import Link from 'next/link'
import { clsx } from 'clsx'
import { useRouter, useSearchParams } from 'next/navigation'
import { errorNotification, materialConfirm, successNotification } from '~/reusable-functions'
import { Icons } from '~/components/icons'

const breadcrumbItems = [{ title: 'lists', link: '/lists' }]

const ListsPage = () => {
	// * 1. Create a page for lists
	// * 2. Create a form to add a contact
	// * 3. Import bulk contact button
	// * 4. Bulk select actions : Export, Delete, Create a new List
	// * 5 . Individual contact actions : Edit, Delete

	const searchParams = useSearchParams()
	const router = useRouter()

	const page = Number(searchParams.get('page') || 1)
	const pageLimit = Number(searchParams.get('limit') || 0) || 10

	const { data: contactListResponse, refetch } = useGetContactLists({
		page: page || 1,
		per_page: pageLimit || 10
	})

	const totalLists = contactListResponse?.paginationMeta?.total || 0
	const pageCount = Math.ceil(totalLists / pageLimit)
	const lists: ContactListSchema[] = contactListResponse?.lists || []

	const deleteContactListMutation = useDeleteListById()
	async function deleteContactList(id: string) {
		try {
			const confirmation = await materialConfirm({
				title: 'Delete Contact List',
				description: 'Are you sure you want to delete this contact list?'
			})

			if (!confirmation) return

			const response = await deleteContactListMutation.mutateAsync({
				id: id
			})

			if (response) {
				successNotification({
					message: 'Contact list deleted successfully'
				})
				await refetch()
			} else {
				errorNotification({
					message: 'Failed to delete contact list'
				})
			}
		} catch (error) {
			console.error('Failed to delete contact list', error)
			errorNotification({
				message: 'Failed to delete contact list'
			})
		}
	}

	return (
		<>
			<div className="flex-1 space-y-4  p-4 pt-6 md:p-8">
				<BreadCrumb items={breadcrumbItems} />

				<div className="flex items-start justify-between">
					<Heading title={`Lists (${lists.length})`} description="Manage lists" />

					<Link
						href={'/lists/new-or-edit'}
						className={clsx(buttonVariants({ variant: 'default' }))}
					>
						<Icons.add className="mr-2 h-4 w-4" /> Add New
					</Link>
				</div>
				<Separator />
				<TableComponent
					searchKey="name"
					pageNo={page}
					columns={ContactListTableColumns}
					totalUsers={totalLists}
					data={lists}
					pageCount={pageCount}
					actions={[
						{
							icon: 'edit',
							label: 'Edit',
							onClick: (contactListId: string) => {
								// redirect to the edit page with id in search param
								router.push(`/lists/new-or-edit?id=${contactListId}`)
							}
						},
						{
							icon: 'trash',
							label: 'Delete',
							onClick: (contactListId: string) => {
								deleteContactList(contactListId).catch(console.error)
							}
						}
					]}
				/>
			</div>
		</>
	)
}

export default ListsPage
