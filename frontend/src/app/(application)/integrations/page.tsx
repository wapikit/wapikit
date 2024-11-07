'use client'

import {
	type IntegrationSchema,
	IntegrationStatusEnum,
	OrderEnum,
	useGetIntegrations
} from 'root/.generated'
import BreadCrumb from '~/components/breadcrumb'
import IntegrationCard from '~/components/integration/card'
import { Heading } from '~/components/ui/heading'
import { ScrollArea } from '~/components/ui/scroll-area'
import { Separator } from '~/components/ui/separator'

const IntegrationsPage = () => {
	const breadcrumbItems = [{ title: 'Integrations', link: '/integrations' }]

	const { data } = useGetIntegrations({
		order: OrderEnum.asc,
		page: 1,
		per_page: 10
	})

	const integrations: IntegrationSchema[] = [
		{
			name: 'Open AI',
			uniqueId: 'open-ai',
			slug: 'open-ai',
			status: IntegrationStatusEnum.Inactive,
			type: 'messaging',
			icon: 'https://media.discordapp.net/attachments/907937769014325288/1304066673174777886/openai-2.png?ex=672e0a3b&is=672cb8bb&hm=9bbebaa23dc9d94927c6331c07049a6e4688e44043ef76dbc35ce0fe6909e4ac&=&format=webp&quality=lossless&width=1214&height=1227',
			description:
				'This integration allows you to integrate with slack for real time notification for you custom support conversation team inbox',
			createdAt: new Date().toISOString(),
			isPremium: false
		},
		{
			name: 'Slack',
			uniqueId: 'slack',
			slug: 'slack',
			status: IntegrationStatusEnum.Inactive,
			icon: 'https://cdn.bfldr.com/5H442O3W/at/pl546j-7le8zk-afym5u/Slack_Mark_Web.png?auto=webp&format=png',
			type: 'custom',
			description:
				'This integrations allows you to integrate with openai for custom support conversation team inbox. Once you provide your API key, the integration will generate responses for incoming messages via the AI model API calls.',
			createdAt: new Date().toISOString(),
			isPremium: true
		},
		{
			name: 'Shopify',
			uniqueId: 'shopify',
			slug: 'shopify',
			status: IntegrationStatusEnum.Inactive,
			type: 'messaging',
			icon: 'https://media.discordapp.net/attachments/907937769014325288/1304066413316804608/images.png?ex=672e09fd&is=672cb87d&hm=34aa3645a39be0457574e62af8323e494172de270cdc02b01744e63d432d315b&=&format=webp&quality=lossless&width=464&height=525',
			description: 'Integrate with shopify to get real time notification',
			createdAt: new Date().toISOString(),
			isPremium: false
		},
		{
			name: 'Wordpress',
			uniqueId: 'wordpress',
			slug: 'wordpress',
			status: IntegrationStatusEnum.Inactive,
			type: 'messaging',
			icon: 'https://media.discordapp.net/attachments/907937769014325288/1304066797141491782/images.png?ex=672e0a58&is=672cb8d8&hm=364081e344e271dc758cf8162768a55872915af53a863c1b41be9b739e576606&=&format=webp&quality=lossless&width=495&height=495',
			description: 'Integrate with wordpress to get real time notification',
			createdAt: new Date().toISOString(),
			isPremium: false
		}
		// {
		// 	name: 'Website Chatbot',
		// 	uniqueId: 'slack',
		// 	slug: 'slack',
		// 	status: IntegrationStatusEnum.Inactive,
		// 	icon: 'https://media.discordapp.net/attachments/1007886641484005427/1269616730968297534/image.png?ex=66b0b639&is=66af64b9&hm=f0c0c74a5f76827f0f1a8f211318a9332119bebe44abe71ecce99c2710c45a72&=&format=webp&quality=lossless&width=1227&height=1227',
		// 	type: 'custom',
		// 	description: '',
		// 	createdAt: new Date().toISOString()
		// },
	]

	console.log(data)

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
