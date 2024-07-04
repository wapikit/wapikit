'use client'
// import { AlertModal } from '~/components/modal/alert-modal'
import { Button } from '~/components/ui/button'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuTrigger
} from '~/components/ui/dropdown-menu'
import { MoreHorizontal } from 'lucide-react'
import { type TableCellActionProps } from '~/types'
import { Icons } from '../icons'

export const CellAction: React.FC<{ actions: TableCellActionProps[] }> = ({ actions }) => {
	// const [loading] = useState(false)
	// const [open, setOpen] = useState(false)

	return (
		<>
			<DropdownMenu modal={false}>
				<DropdownMenuTrigger asChild>
					<Button variant="ghost" className="h-8 w-8 p-0">
						<span className="sr-only">Open menu</span>
						<MoreHorizontal className="h-4 w-4 " />
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
									action.onClick()
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
