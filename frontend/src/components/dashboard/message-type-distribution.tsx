import { BarChart } from '@tremor/react'

const data = [
	{
		name: 'Text',
		sent: Math.floor(Math.random() * 5000) + 1000,
		received: Math.floor(Math.random() * 5000) + 1000
	},
	{
		name: 'Video',
		sent: Math.floor(Math.random() * 5000) + 1000,
		received: Math.floor(Math.random() * 5000) + 1000
	},
	{
		name: 'Audio',
		sent: Math.floor(Math.random() * 5000) + 1000,
		received: Math.floor(Math.random() * 5000) + 1000
	},
	{
		name: 'Document',
		sent: Math.floor(Math.random() * 5000) + 1000,
		received: Math.floor(Math.random() * 5000) + 1000
	},
	{
		name: 'Location',
		sent: Math.floor(Math.random() * 5000) + 1000,
		received: Math.floor(Math.random() * 5000) + 1000
	},
	{
		name: 'Button Interactions',
		sent: Math.floor(Math.random() * 5000) + 1000,
		received: Math.floor(Math.random() * 5000) + 1000
	},
	{
		name: 'Image',
		sent: Math.floor(Math.random() * 5000) + 1000,
		received: Math.floor(Math.random() * 5000) + 1000
	},
	{
		name: 'Contact',
		sent: Math.floor(Math.random() * 5000) + 2000,
		received: Math.floor(Math.random() * 5000) + 2000
	},
	{
		name: 'List',
		sent: Math.floor(Math.random() * 5000) + 2000,
		received: Math.floor(Math.random() * 5000) + 2000
	},
	{
		name: 'List Reply',
		sent: Math.floor(Math.random() * 5000) + 2000,
		received: Math.floor(Math.random() * 5000) + 2000
	}
]

export function MessageTypeBifurcation() {
	return (
		<div className="w-full rounded-lg">
			<BarChart
				data={data}
				index="name"
				categories={['sent', 'received']}
				colors={['green', 'indigo']}
				yAxisWidth={48}
				showAnimation
			/>
		</div>
	)
}
