'use client'

import { LineChart } from '@tremor/react'
import React from 'react'
import { type LinkClicksGraphDataPointSchema } from 'root/.generated'

export const LinkClicks: React.FC<{ data: LinkClicksGraphDataPointSchema[] }> = ({ data }) => {
	return (
		<div className="h-[375px] w-full rounded-lg">
			<LineChart
				className="mt-20 text-xs"
				data={data || []}
				index="label"
				categories={['count']}
				colors={['green']}
				showLegend={false}
				showAnimation
				showTooltip={true}
				curveType="natural"
			/>
		</div>
	)
}
