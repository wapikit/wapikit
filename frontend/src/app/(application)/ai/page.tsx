'use client'

import { AiChat } from '~/components/ai/ai-chat'
import { useAiChatStore } from '~/store/ai-chat-store'

const AiChatBox = () => {
	const { chats } = useAiChatStore()
	const chat = chats[0]
	if (!chat) {
		return null
	}
	return <AiChat chat={chat} />
}

export default AiChatBox
