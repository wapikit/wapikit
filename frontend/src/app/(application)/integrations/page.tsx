'use client'

import { OrderEnum, useGetIntegrations } from 'root/.generated'
import { ScrollArea } from '~/components/ui/scroll-area'

const IntegrationsPage = () => {
	const { data } = useGetIntegrations({
		order: OrderEnum.asc,
		page: 1,
		per_page: 10
	})

	console.log(data)

	return <ScrollArea className="h-full"></ScrollArea>
}

export default IntegrationsPage
