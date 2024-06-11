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

export const CellAction: React.FC<{ actions: TableCellActionProps[] }> = ({ actions }) => {
	// const [loading] = useState(false)
	// const [open, setOpen] = useState(false)

	return (
		<>
			{/* <AlertModal
				isOpen={open}
				onClose={() => setOpen(false)}
				onConfirm={onConfirm}
				loading={loading}
			/> */}
			<DropdownMenu modal={false}>
				<DropdownMenuTrigger asChild>
					<Button variant="ghost" className="h-8 w-8 p-0">
						<span className="sr-only">Open menu</span>
						<MoreHorizontal className="h-4 w-4" />
					</Button>
				</DropdownMenuTrigger>
				<DropdownMenuContent align="end">
					<DropdownMenuLabel>Actions</DropdownMenuLabel>

					{actions.map((action, index) => {
						return (
							<DropdownMenuItem
								key={index}
								onClick={() => {
									action.onClick()
								}}
							>
								{action.icon} {action.label}
							</DropdownMenuItem>
						)
					})}
				</DropdownMenuContent>
			</DropdownMenu>
		</>
	)
}
