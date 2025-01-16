'use client'

import { type ChatBotStateEnum } from '~/types'
import { Badge } from '../ui/badge'
import { Separator } from '../ui/separator'

const ChatHeader = ({
	chatTitle,
	chatBotState
}: {
	chatTitle: string
	chatBotState: ChatBotStateEnum
}) => {
	return (
		<header className="sticky top-0 flex flex-col items-center gap-2 rounded-md bg-opacity-35  px-2 py-1.5 md:px-2">
			<div className="flex w-full items-center justify-start gap-2 text-sm">
				{chatTitle}
				<Badge>{chatBotState}</Badge>
			</div>
			<Separator className="my-2" />
		</header>
	)
}

export default ChatHeader
