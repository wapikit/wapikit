'use client'

import { useSearchParams } from 'next/navigation'
import { useDeleteListById, useGetListById } from 'root/.generated'
import BreadCrumb from '~/components/breadcrumb'
import DocumentationPitch from '~/components/forms/documentation-pitch'
import NewContactListForm from '~/components/forms/new-contact-list-form'
import { Icons } from '~/components/icons'
import { Button } from '~/components/ui/button'
import { Heading } from '~/components/ui/heading'
import { ScrollArea } from '~/components/ui/scroll-area'
import { Separator } from '~/components/ui/separator'
import { errorNotification, materialConfirm, successNotification } from '~/reusable-functions'

const CreateNewContactListPage = () => {
	const breadcrumbItems = [
		{ title: 'Lists', link: '/lists' },
		{ title: 'Create', link: '/lists/new-or-edit' }
	]

	const searchParams = useSearchParams()
	const listId = searchParams.get('id')
	const deleteContactListMutation = useDeleteListById()

	const listResponse = useGetListById(listId || '', {
		query: {
			enabled: !!listId
		}
	})

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
		<ScrollArea className="h-full">
			<div className="flex-1 space-y-4  p-4 pt-6 md:p-8">
				<BreadCrumb items={breadcrumbItems} />
				<div className="flex items-start justify-between">
					<Heading
						title={listId ? 'Edit Contact List' : `Create New Contact List`}
						description=""
					/>
					{listId && (
						<Button
							variant="destructive"
							className="flex items-center gap-2"
							onClick={() => {
								deleteContactList(listId).catch(error => console.error(error))
							}}
						>
							<Icons.trash className="h-4 w-4" />
							Delete List
						</Button>
					)}
				</div>
				<Separator />

				<div className="flex flex-row gap-10">
					<NewContactListForm initialData={listResponse.data?.list || null} />
					<DocumentationPitch type="lists" />
				</div>
			</div>
		</ScrollArea>
	)
}

export default CreateNewContactListPage
