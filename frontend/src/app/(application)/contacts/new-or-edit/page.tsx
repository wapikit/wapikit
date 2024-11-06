'use client'

import { useSearchParams } from 'next/navigation'
import { useGetContactById } from 'root/.generated'
import BreadCrumb from '~/components/breadcrumb'
import DocumentationPitch from '~/components/forms/documentation-pitch'
import NewContactForm from '~/components/forms/new-contact-form'
import { Heading } from '~/components/ui/heading'
import { Separator } from '~/components/ui/separator'

const CreateNewContactPage = () => {
	const breadcrumbItems = [
		{ title: 'Contacts', link: '/contacts' },
		{ title: 'Create', link: '/contacts/new-or-edit' }
	]

	const searchParams = useSearchParams()
	const contactId = searchParams.get('id')

	const contactResponse = useGetContactById(contactId || '', {
		query: {
			enabled: !!contactId
		}
	})

	return (
		<>
			<div className="flex-1 space-y-4  p-4 pt-6 md:p-8">
				<BreadCrumb items={breadcrumbItems} />
				<div className="flex items-start justify-between">
					<Heading title={`Create New Contact`} description="" />
				</div>
				<Separator />

				<div className="flex flex-row gap-10">
					<NewContactForm initialData={contactResponse.data?.contact || null} />
					<DocumentationPitch type="contact" />
				</div>
			</div>
		</>
	)
}

export default CreateNewContactPage
