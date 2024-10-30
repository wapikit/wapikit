'use client'

import { LineChart } from '@tremor/react'
import React from 'react'
import { type MessageAnalyticGraphDataPointSchema } from 'root/.generated'

export const MessageAggregateAnalytics: React.FC<{
	data: MessageAnalyticGraphDataPointSchema[]
}> = ({ data }) => {
	return (
		<div className="h-[375px] w-full rounded-lg">
			<LineChart
				className="mt-20"
				data={data || []}
				index="label"
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
