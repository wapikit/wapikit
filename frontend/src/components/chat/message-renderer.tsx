import { clsx } from 'clsx'
import dayjs from 'dayjs'
import { MessageDirectionEnum, type MessageSchema } from 'root/.generated'
import { ChevronDown } from 'lucide-react'
import {
	DropdownMenu,
	DropdownMenuTrigger,
	DropdownMenuContent,
	DropdownMenuItem
} from '../ui/dropdown-menu'
import { Icons } from '../icons'
import { successNotification } from '~/reusable-functions'
import { useCopyToClipboard } from 'usehooks-ts'

// ! TODO complete this component, right now only supports text messages

const MessageRenderer: React.FC<{ message: MessageSchema; isActionsEnabled: boolean }> = ({
	message,
	isActionsEnabled
}) => {
	const copyToClipboard = useCopyToClipboard()[1]

	const messageActions: {
		label: string
		icon: keyof typeof Icons
		onClick?: () => void
	}[] = [
		{
			label: 'Delete',
			icon: 'trash',
			// ! TODO: implement delete message
			onClick: () => {}
		},
		{
			label: 'Reply',
			icon: 'reply',
			// ! TODO: implement reply message
			onClick: () => {}
		},
		{
			label: 'Copy',
			icon: 'clipboard',
			onClick: () => {
				copyToClipboard((message.messageData || '') as string).catch(error =>
					console.error(error)
				)
				successNotification({
					message: 'Copied'
				})
			}
		}
	]

	return (
		<div
			className={clsx(
				'flex  w-fit max-w-md gap-1  rounded-md p-1 px-3',
				message.direction === MessageDirectionEnum.InBound
					? 'mr-auto bg-white dark:bg-[#202c33]'
					: 'ml-auto bg-primary  text-primary-foreground dark:bg-[#005c4b]'
			)}
		>
			{message.message_type === 'Text' ? (
				<p className="text-wrap text-sm">{message.messageData?.text as any}</p>
			) : null}

			<div className="flex flex-col items-center  justify-end gap-1">
				{isActionsEnabled ? (
					<div className="ml-auto">
						<DropdownMenu modal={false}>
							<DropdownMenuTrigger asChild>
								<ChevronDown
									className={clsx(
										'text-bold h-5 w-5',
										message.direction === MessageDirectionEnum.InBound
											? ''
											: ' text-primary-foreground'
									)}
								/>
							</DropdownMenuTrigger>
							<DropdownMenuContent align="end" side="right">
								{messageActions.map((action, index) => {
									const Icon = Icons[action.icon]
									return (
										<DropdownMenuItem
											key={index}
											onClick={() => {
												if (action.onClick) {
													action.onClick()
												}
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
					</div>
				) : null}

				{message.createdAt ? (
					<span
						className={clsx(
							'ml-auto text-[10px]',
							message.direction === MessageDirectionEnum.InBound
								? ''
								: 'text-primary-foreground'
						)}
					>
						{dayjs(message.createdAt).format('hh:mm A')}
					</span>
				) : null}
			</div>
		</div>
	)
}

export default MessageRenderer
