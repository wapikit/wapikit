'use client'

import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuTrigger
} from '~/components/ui/dropdown-menu'
import { type TableCellActionProps } from '~/types'
import { Icons } from '../icons'
import { MoreVerticalIcon } from 'lucide-react'

export const CellAction: React.FC<{ actions: TableCellActionProps[]; data: any }> = ({
	actions,
	data
}) => {
	console.log('data', data)

	return (
		<>
			<DropdownMenu modal={false}>
				<DropdownMenuTrigger asChild>
					<MoreVerticalIcon className="h-4 w-4 text-secondary-foreground" />
				</DropdownMenuTrigger>
				<DropdownMenuContent align="end">
					<DropdownMenuLabel>Actions</DropdownMenuLabel>
					{actions.map((action, index) => {
						const Icon = Icons[action.icon]
						return (
							<DropdownMenuItem
								key={index}
								onClick={() => {
									// @ts-ignore
									action.onClick(data)
								}}
								className="flex flex-row items-center gap-2"
								disabled={action.disabled || false}
							>
								<Icon className="size-4" />
								{action.label}
							</DropdownMenuItem>
						)
					})}
				</DropdownMenuContent>
			</DropdownMenu>
		</>
	)
}
