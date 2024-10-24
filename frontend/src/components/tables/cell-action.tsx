'use client'

import { Button } from '~/components/ui/button'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuTrigger
} from '~/components/ui/dropdown-menu'
import { type TableCellActionProps } from '~/types'
import { Icons } from '../icons'

export const CellAction: React.FC<{ actions: TableCellActionProps[]; data: any }> = ({
	actions,
	data
}) => {
	console.log('data', data)
	// const [loading] = useState(false)
	// const [open, setOpen] = useState(false)

	const MoreIcon = Icons['ellipsis']

	return (
		<>
			<DropdownMenu modal={false}>
				<DropdownMenuTrigger asChild>
					<Button variant="ghost" className="h-8 w-8 p-0 text-foreground">
						<span className="sr-only">Open menu</span>
						More
						<MoreIcon className="h-4 w-4 text-green-500" />
					</Button>
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
