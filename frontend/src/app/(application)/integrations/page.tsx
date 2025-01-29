'use client'

import { type IntegrationSchema, IntegrationStatusEnum } from 'root/.generated'
import BreadCrumb from '~/components/breadcrumb'
import IntegrationCard from '~/components/integration/card'
import { Heading } from '~/components/ui/heading'
import { ScrollArea } from '~/components/ui/scroll-area'
import { Separator } from '~/components/ui/separator'

const IntegrationsPage = () => {
	const breadcrumbItems = [{ title: 'Integrations', link: '/integrations' }]

	// const { data } = useGetIntegrations({
	// 	order: OrderEnum.asc,
	// 	page: 1,
	// 	per_page: 10
	// })

	const integrations: IntegrationSchema[] = [
		{
			name: 'HubSpot',
			uniqueId: 'hubspot',
			slug: 'hubspot',
			status: IntegrationStatusEnum.Inactive,
			type: 'messaging',
			icon: '/assets/integrations/hubspot.svg',
			description:
				'This integration allows you to export campaign contacts or leads to HubSpot',
			createdAt: new Date().toISOString(),
			isPremium: false
		},
		{
			name: 'Salesforce',
			uniqueId: 'sales-force',
			slug: 'open-ai',
			status: IntegrationStatusEnum.Inactive,
			type: 'messaging',
			icon: '/assets/integrations/salesforce.svg',
			description:
				'This integration allows you to export campaign contacts or leads to HubSpot',
			createdAt: new Date().toISOString(),
			isPremium: false
		},
		{
			name: 'Shopify',
			uniqueId: 'shopify',
			slug: 'shopify',
			status: IntegrationStatusEnum.Inactive,
			type: 'messaging',
			icon: '/assets/integrations/shopify.svg',
			description: 'Sync customer data and send WhatsApp notifications for order updates.',
			createdAt: new Date().toISOString(),
			isPremium: false
		},
		{
			name: 'WooCommerce',
			uniqueId: 'woocommerce',
			slug: 'woocommerce',
			status: IntegrationStatusEnum.Inactive,
			type: 'messaging',
			icon: '/assets/integrations/woocommerce.svg',
			description: 'Send abandoned cart reminders or delivery updates via WhatsApp.',
			createdAt: new Date().toISOString(),
			isPremium: true
		},
		{
			name: 'Google Sheets',
			uniqueId: 'google-sheets',
			slug: 'google-sheets',
			status: IntegrationStatusEnum.Inactive,
			type: 'messaging',
			icon: '/assets/integrations/sheets.svg',
			description: 'Export conversation or campaign data to spreadsheets for analysis.',
			createdAt: new Date().toISOString(),
			isPremium: true
		},
		{
			name: 'Notion',
			uniqueId: 'notion',
			slug: 'notion',
			status: IntegrationStatusEnum.Inactive,
			type: 'messaging',
			icon: '/assets/integrations/notion.svg',
			description:
				'Automatically log campaign insights or conversation summaries into shared team documents.',
			createdAt: new Date().toISOString(),
			isPremium: false
		},
		{
			name: 'Linear',
			uniqueId: 'linear',
			slug: 'linear',
			status: IntegrationStatusEnum.Inactive,
			type: 'messaging',
			icon: '/assets/integrations/linear.svg',
			description: 'Create issues right from your whatsapp communications.',
			createdAt: new Date().toISOString(),
			isPremium: false
		},
		{
			name: 'Razorpay',
			uniqueId: 'razorpay',
			slug: 'razorpay',
			status: IntegrationStatusEnum.Inactive,
			type: 'messaging',
			icon: '/assets/integrations/razorpay.svg',
			description:
				'Send payment links to customers via WhatsApp, manage your ecommerce store payments and more.',
			createdAt: new Date().toISOString(),
			isPremium: false
		},
		{
			name: 'Zapier',
			uniqueId: 'zapier',
			slug: 'zapier',
			status: IntegrationStatusEnum.Inactive,
			type: 'messaging',
			icon: '/assets/integrations/zapier.svg',
			description: 'Allow users to create custom workflows with thousands of other apps.',
			createdAt: new Date().toISOString(),
			isPremium: false
		},
		{
			name: 'Make',
			uniqueId: 'make',
			slug: 'make',
			status: IntegrationStatusEnum.Inactive,
			type: 'messaging',
			icon: '/assets/integrations/make.svg',
			description: 'Similar to Zapier, enables users to connect WapiKit with other tools.',
			createdAt: new Date().toISOString(),
			isPremium: false
		}
	]

	return (
		<ScrollArea className="h-full">
			<div className="flex-1 space-y-4  p-4 pt-6 md:p-8">
				<BreadCrumb items={breadcrumbItems} />
				<div className="flex items-start justify-between">
					<Heading title={`Integrations`} description="" />
				</div>
				<Separator />
			</div>
			<section className="mr-auto grid max-w-6xl grid-cols-3 flex-wrap gap-5 pl-8">
				{integrations.map((integration, index) => {
					return (
						<IntegrationCard
							key={index}
							icon={integration.icon}
							name={integration.name}
							slug={integration.slug}
							description={integration.description}
						/>
					)
				})}
			</section>
		</ScrollArea>
	)
}

export default IntegrationsPage
