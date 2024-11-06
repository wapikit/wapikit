'use client'

import { LineChart } from '@tremor/react'

const data = [
	{
		name: `${new Date().getDate() - 10} ${new Date().toString().split(' ')[1]}`,
		total: Math.floor(Math.random() * 5000) + 1000
	},
	{
		name: `${new Date().getDate() - 9} ${new Date().toString().split(' ')[1]}`,
		total: Math.floor(Math.random() * 5000) + 1000
	},
	{
		name: `${new Date().getDate() - 8} ${new Date().toString().split(' ')[1]}`,
		total: Math.floor(Math.random() * 5000) + 1000
	},
	{
		name: `${new Date().getDate() - 7} ${new Date().toString().split(' ')[1]}`,
		total: Math.floor(Math.random() * 5000) + 1000
	},
	{
		name: `${new Date().getDate() - 6} ${new Date().toString().split(' ')[1]}`,
		total: Math.floor(Math.random() * 5000) + 1000
	},
	{
		name: `${new Date().getDate() - 5} ${new Date().toString().split(' ')[1]}`,
		total: Math.floor(Math.random() * 5000) + 1000
	},
	{
		name: `${new Date().getDate() - 4} ${new Date().toString().split(' ')[1]}`,
		total: Math.floor(Math.random() * 5000) + 1000
	},
	{
		name: `${new Date().getDate() - 3} ${new Date().toString().split(' ')[1]}`,
		total: Math.floor(Math.random() * 5000) + 1000
	},
	{
		name: `${new Date().getDate() - 2} ${new Date().toString().split(' ')[1]}`,
		total: Math.floor(Math.random() * 5000) + 1000
	},
	{
		name: `${new Date().getDate() - 1} ${new Date().toString().split(' ')[1]}`,
		total: Math.floor(Math.random() * 5000) + 1000
	},
	{
		name: `${new Date().getDate()} ${new Date().toString().split(' ')[1]}`,
		total: Math.floor(Math.random() * 5000) + 1000
	}
]

export function Replies() {
	return (
		<div className="h-[375px] w-full rounded-lg">
			<LineChart
				className="mt-14 text-xs"
				data={data || []}
				index="name"
				categories={['total']}
				colors={['green']}
				showLegend={false}
				showAnimation
				style={{ stroke: 'green' }}
				showTooltip={true}
				curveType="natural"
				unselectable="on"
			/>
		</div>
	)
}
