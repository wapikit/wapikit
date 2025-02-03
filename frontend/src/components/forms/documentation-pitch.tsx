import Link from 'next/link'
import React from 'react'
import { Icons } from '~/components/icons'
import { Button } from '~/components/ui/button'

const data = [
	{
		slug: 'campaign',
		title: 'Campaigns',
		description: 'Checkout the documentation for campaigns, how to create and manage them.',
		Icon: Icons.rocket,
		ctaText: 'Check Docs',
		ctaUrl: 'https://docs.wapikit.com/guide/manage-campaigns'
	},
	{
		slug: 'contact',
		title: 'Contacts',
		description: 'Checkout the documentation for contacts, how to create and manage them.',
		Icon: Icons.profile,
		ctaText: 'Check Docs',
		ctaUrl: 'https://docs.wapikit.com/guide/manage-contacts'
	},
	{
		slug: 'lists',
		title: 'Lists',
		description: 'Checkout the documentation for lists, how to create and manage them.',
		Icon: Icons.laptop,
		ctaText: 'Check Docs',
		ctaUrl: 'https://docs.wapikit.com/guide/manage-lists'
	},
	{
		slug: 'api-key',
		title: 'Api Key Usage',
		description: 'Checkout the documentation for api key usage, how to create and manage them.',
		Icon: Icons.login,
		ctaText: 'Check Docs',
		ctaUrl: 'https://docs.wapikit.com/guide/manage-api-keys'
	}
] as const

type DocumentationPitchType = (typeof data)[number]['slug']

const DocumentationPitch: React.FC<{ type: DocumentationPitchType }> = ({ type }) => {
	const dataToUse = data.find(item => item.slug === type)

	if (!dataToUse) {
		return null
	} else {
		const { Icon, ctaText, ctaUrl, description, slug, title } = dataToUse

		return (
			<div
				className="documentation-pitch group mt-8 flex h-fit max-w-md flex-col gap-4 rounded-lg border p-4 hover:border-primary"
				key={slug}
			>
				<div className="flex items-center gap-2">
					<Icon className={`size-5 group-hover:text-primary`} />
					<h2 className="text-lg font-semibold ">{title}</h2>
				</div>
				<p className="">{description}</p>
				<Link href={ctaUrl}>
					<Button>{ctaText}</Button>
				</Link>
			</div>
		)
	}
}

export default DocumentationPitch
