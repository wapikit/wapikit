import { Badge } from '../ui/badge'
import { Modal } from '../ui/modal'

const RenderBadgedListModel = (params: {
	isOpen: boolean
	setIsOpen: (value: boolean) => void
	content: string[]
	title: string
	description?: string
}) => {
	const { isOpen, setIsOpen, content, title, description } = params

	return (
		<Modal
			title={title}
			description={description || ''}
			isOpen={isOpen}
			onClose={() => {
				setIsOpen(false)
			}}
		>
			<div className="flex flex-wrap items-center justify-center gap-0.5 truncate">
				{content.map((label, index) => {
					if (index > 2) {
						return null
					}
					return <Badge key={index}>{label}</Badge>
				})}
			</div>
		</Modal>
	)
}

export default RenderBadgedListModel
