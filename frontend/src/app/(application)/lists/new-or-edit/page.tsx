'use client'

import { useSearchParams } from 'next/navigation'
import { useGetCampaignById, useGetListById } from 'root/.generated'
import BreadCrumb from '~/components/breadcrumb'
import DocumentationPitch from '~/components/forms/documentation-pitch'
import NewContactListForm from '~/components/forms/new-contact-list-form'
import { Heading } from '~/components/ui/heading'
import { ScrollArea } from '~/components/ui/scroll-area'
import { Separator } from '~/components/ui/separator'

const CreateNewContactListPage = () => {
	const breadcrumbItems = [
		{ title: 'Lists', link: '/lists' },
		{ title: 'Create', link: '/lists/new-or-edit' }
	]

	const searchParams = useSearchParams()
	const listId = searchParams.get('id')

	const listResponse = useGetListById(listId || '', {
		query: {
			enabled: !!listId
		}
	})

	return (
		<ScrollArea className="h-full">
			<div className="flex-1 space-y-4  p-4 pt-6 md:p-8">
				<BreadCrumb items={breadcrumbItems} />
				<div className="flex items-start justify-between">
					<Heading title={`Create New Contact List`} description="" />
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
