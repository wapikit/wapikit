'use client'

import { LineChart } from '@tremor/react'

const data = [
	{
		name: `${new Date().getDate() - 10} ${new Date().toString().split(' ')[1]}`,
		sent: Math.floor(Math.random() * 5000) + 1000,
		replied: Math.floor(Math.random() * 5000) + 1000,
		read: Math.floor(Math.random() * 5000) + 1000
	},
	{
		name: `${new Date().getDate() - 9} ${new Date().toString().split(' ')[1]}`,
		sent: Math.floor(Math.random() * 5000) + 1000,
		replied: Math.floor(Math.random() * 5000) + 1000,
		read: Math.floor(Math.random() * 5000) + 1000
	},
	{
		name: `${new Date().getDate() - 8} ${new Date().toString().split(' ')[1]}`,
		sent: Math.floor(Math.random() * 5000) + 1000,
		replied: Math.floor(Math.random() * 5000) + 1000,
		read: Math.floor(Math.random() * 5000) + 1000
	},
	{
		name: `${new Date().getDate() - 7} ${new Date().toString().split(' ')[1]}`,
		sent: Math.floor(Math.random() * 5000) + 1000,
		replied: Math.floor(Math.random() * 5000) + 1000,
		read: Math.floor(Math.random() * 5000) + 1000
	},
	{
		name: `${new Date().getDate() - 6} ${new Date().toString().split(' ')[1]}`,
		sent: Math.floor(Math.random() * 5000) + 1000,
		replied: Math.floor(Math.random() * 5000) + 1000,
		read: Math.floor(Math.random() * 5000) + 1000
	},
	{
		name: `${new Date().getDate() - 5} ${new Date().toString().split(' ')[1]}`,
		sent: Math.floor(Math.random() * 5000) + 1000,
		replied: Math.floor(Math.random() * 5000) + 1000,
		read: Math.floor(Math.random() * 5000) + 1000
	},
	{
		name: `${new Date().getDate() - 4} ${new Date().toString().split(' ')[1]}`,
		sent: Math.floor(Math.random() * 5000) + 1000,
		replied: Math.floor(Math.random() * 5000) + 1000,
		read: Math.floor(Math.random() * 5000) + 1000
	},
	{
		name: `${new Date().getDate() - 3} ${new Date().toString().split(' ')[1]}`,
		sent: Math.floor(Math.random() * 5000) + 1000,
		replied: Math.floor(Math.random() * 5000) + 1000,
		read: Math.floor(Math.random() * 5000) + 1000
	},
	{
		name: `${new Date().getDate() - 2} ${new Date().toString().split(' ')[1]}`,
		sent: Math.floor(Math.random() * 5000) + 1000,
		replied: Math.floor(Math.random() * 5000) + 1000,
		read: Math.floor(Math.random() * 5000) + 1000
	},
	{
		name: `${new Date().getDate() - 1} ${new Date().toString().split(' ')[1]}`,
		sent: Math.floor(Math.random() * 5000) + 1000,
		replied: Math.floor(Math.random() * 5000) + 1000,
		read: Math.floor(Math.random() * 5000) + 1000
	},
	{
		name: `${new Date().getDate()} ${new Date().toString().split(' ')[1]}`,
		sent: Math.floor(Math.random() * 5000) + 1000,
		replied: Math.floor(Math.random() * 5000) + 1000,
		read: Math.floor(Math.random() * 5000) + 1000
	}
]

export function MessageAggregateAnalytics() {
	return (
		<div className="h-[375px] w-full rounded-lg">
			<LineChart
				className="mt-20"
				data={data || []}
				index="name"
				categories={['sent', 'read', 'replied']}
				colors={['blue', 'yellow', 'green']}
				showLegend={false}
				showAnimation
				showTooltip={true}
				curveType="natural"
				unselectable="on"
			/>
		</div>
	)
}
