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
			name: 'Slack',
			uniqueId: 'slack',
			slug: 'slack',
			status: IntegrationStatusEnum.Inactive,
			type: 'messaging',
			icon: 'https://media.discordapp.net/attachments/1007886641484005427/1269616730968297534/image.png?ex=66b0b639&is=66af64b9&hm=f0c0c74a5f76827f0f1a8f211318a9332119bebe44abe71ecce99c2710c45a72&=&format=webp&quality=lossless&width=1227&height=1227',
			description:
				'This integration allows you to integrate with slack for real time notification for you custom support conversation team inbox',
			createdAt: new Date().toISOString(),
			isPremium: false
		},
		{
			name: 'Discord',
			uniqueId: 'discord',
			slug: 'discord',
			status: IntegrationStatusEnum.Inactive,
			type: 'messaging',
			icon: 'https://media.discordapp.net/attachments/1007886641484005427/1269619779476652155/Frame_39_1.png?ex=66b0b90f&is=66af678f&hm=465f06327f7cb398033f81f390dd01090c154ea91e5b5b4211291ec25366f073&=&format=webp&quality=lossless&width=895&height=895',
			description:
				'This integration allows you to integrate with slack for real time notification for you custom support conversation team inbox',
			createdAt: new Date().toISOString(),
			isPremium: false
		},
		{
			name: 'Shopify',
			uniqueId: 'slack',
			slug: 'slack',
			status: IntegrationStatusEnum.Inactive,
			type: 'messaging',
			icon: 'https://media.discordapp.net/attachments/1007886641484005427/1269620042702651467/image_10.png?ex=66b0b94e&is=66af67ce&hm=ffdfd5b27f8cc0d12a258c6109a02b95cb2fbaf2fd0392d517ae39a343e1a57c&=&format=webp&quality=lossless&width=770&height=770',
			description: 'Integrate with shopify to get real time notification',
			createdAt: new Date().toISOString(),
			isPremium: false
		},
		// {
		// 	name: 'Wordpress',
		// 	uniqueId: 'slack',
		// 	slug: 'slack',
		// 	status: IntegrationStatusEnum.Inactive,
		// 	type: 'messaging',
		// 	icon: 'https://media.discordapp.net/attachments/1007886641484005427/1269616730968297534/image.png?ex=66b0b639&is=66af64b9&hm=f0c0c74a5f76827f0f1a8f211318a9332119bebe44abe71ecce99c2710c45a72&=&format=webp&quality=lossless&width=1227&height=1227',
		// 	description: 'Integrate with wordpress to get real time notification',
		// 	createdAt: new Date().toISOString()
		// },
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
		{
			name: 'OpenAI',
			uniqueId: 'slack',
			slug: 'slack',
			status: IntegrationStatusEnum.Inactive,
			icon: 'https://media.discordapp.net/attachments/1007886641484005427/1269619425523535942/Frame.png?ex=66b0b8bb&is=66af673b&hm=01a3b5564587a3571000f739066d672ec61adc43e63ffcf2e590729f31dbe1a1&=&format=webp&quality=lossless&width=895&height=895',
			type: 'custom',
			description:
				'This integrations allows you to integrate with openai for custom support conversation team inbox. Once you provide your API key, the integration will generate responses for incoming messages via the AI model API calls.',
			createdAt: new Date().toISOString(),
			isPremium: true
		}
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
