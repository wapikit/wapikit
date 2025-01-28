'use client'

import { Sheet, SheetContent } from '../ui/sheet'
import { AiChat } from './ai-chat'
import { useAiChatStore } from '~/store/ai-chat-store'

const AiChatBox = () => {
	const { isOpen, writeProperty, chats } = useAiChatStore()
	const chat = chats[0]
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
			<SheetContent
				className="!p-4 sm:!max-w-[750px]"
				onInteractOutside={event => event.preventDefault()}
			>
				<AiChat chat={chat} />
			</SheetContent>
		</Sheet>
	)
}

export default AiChatBox
