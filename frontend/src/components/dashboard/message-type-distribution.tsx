import { BarChart } from '@tremor/react'
import { type MessageTypeDistributionGraphDataPointSchema } from 'root/.generated'

export const MessageTypeBifurcation: React.FC<{
	data: MessageTypeDistributionGraphDataPointSchema[]
}> = ({ data }) => {
	return (
		<div className="w-full rounded-lg">
			<BarChart
				data={data}
				index="label"
				categories={['sent', 'received']}
				colors={['green', 'indigo']}
				yAxisWidth={48}
				showAnimation
			/>
		</div>
	)
}
