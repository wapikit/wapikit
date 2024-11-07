'use client'

import { DonutChart, Legend } from '@tremor/react'

const data = [
	{
		name: 'In Progress',
		value: Math.floor(Math.random() * 5000) + 1000
	},
	{
		name: 'Unread',
		value: Math.floor(Math.random() * 5000) + 1000
	},
	{
		name: 'Unassigned',
		value: Math.floor(Math.random() * 5000) + 1000
	}
]

export function ConversationStatusChart() {
	return (
		<div className="my-auto w-full rounded-lg">
			<Legend
				className="justify-end"
				categories={['In Progress', 'Unread', 'Unassigned']}
				colors={['green', 'indigo', 'red']}
			></Legend>
			<DonutChart
				data={data}
				index="name"
				colors={['green', 'indigo', 'red']}
				showLabel
				showAnimation
			/>
		</div>
	)
}
