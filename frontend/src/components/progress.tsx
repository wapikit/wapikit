import * as React from 'react'

interface ProgressProps {
	value: number
}

const Progress = ({ value }: ProgressProps) => (
	<div className="h-2 rounded-full bg-gray-200">
		<div
			className="h-full rounded-full bg-primary transition-all duration-300"
			style={{ width: `${value}%` }}
		/>
	</div>
)

export default Progress
