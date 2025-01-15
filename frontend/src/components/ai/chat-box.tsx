'use client'

import { Sheet, SheetContent, SheetHeader, SheetTitle } from '../ui/sheet'
import { AiChat } from './chat'
import { convertToUIMessages } from '~/utils/ai-utils'
import { useAiChatStore } from '~/store/ai-chat-store'
import { DataStreamHandler } from './data-stream-handler'

const AiChatBox = () => {
	const { isOpen, writeProperty, chats } = useAiChatStore()
	const chat = chats[0]
	console.log({ chat })
	if (!chat) {
		return null
	}
	return (
		<Sheet
			open={isOpen}
			onOpenChange={isOpen => {
				writeProperty({
					isOpen
				})
			}}
		>
			<SheetHeader>
				<SheetTitle>AI Chat</SheetTitle>
			</SheetHeader>

			<SheetContent
				className="!p-4 sm:!max-w-[750px]"
				onInteractOutside={event => event.preventDefault()}
			>
				<AiChat chat={chat} initialMessages={convertToUIMessages([])} />
				<DataStreamHandler id={chat.uniqueId} />
			</SheetContent>
		</Sheet>
	)
}

export default AiChatBox
