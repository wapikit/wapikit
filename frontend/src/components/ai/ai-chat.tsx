'use client'

import type { Attachment } from 'ai'
import { useEffect, useState } from 'react'
import { AiChatInput } from './ai-chat-input'
import { Messages } from './messages'
import { type AiChatSchema, useGetAiChatMessages, useGetAiChatMessageVotes } from 'root/.generated'
import useChat from '~/hooks/use-chat'
import { ChatBotStateEnum } from '~/types'
import { useAiChatStore } from '~/store/ai-chat-store'

export function AiChat({ chat }: { chat: AiChatSchema }) {
	const { handleSubmit, chatBotState, selectSuggestedAction, currentMessageIdInStream } = useChat(
		{
			chatId: chat.uniqueId
		}
	)

	const { writeProperty } = useAiChatStore()

	// ! TODO: handle pagination here
	const { data: votes } = useGetAiChatMessageVotes(chat.uniqueId, {
		page: 1,
		per_page: 100
	})

	const { data: messages } = useGetAiChatMessages(chat.uniqueId, {
		page: 1,
		per_page: 50
	})

	useEffect(() => {
		writeProperty({
			currentChatMessages: messages?.messages || []
		})
	}, [messages?.messages, writeProperty])

	const [attachments, setAttachments] = useState<Array<Attachment>>([])

	return (
		<div className="flex h-dvh w-full min-w-0 flex-col bg-background px-4 md:px-6">
			{/* <ChatHeader chatTitle={chat.title} chatBotState={chatBotState} /> */}
			<Messages
				chatId={chat.uniqueId}
				isLoading={chatBotState === ChatBotStateEnum.Streaming}
				votes={votes?.votes}
				isReadonly={false}
				chatBotState={chatBotState}
				currentMessageIdInStream={currentMessageIdInStream.current}
			/>

			<form className="bottom-0 mx-auto flex w-full gap-2 bg-background pb-20">
				<AiChatInput
					chatId={chat.uniqueId}
					handleSubmit={handleSubmit}
					isLoading={chatBotState === ChatBotStateEnum.Streaming}
					selectSuggestedAction={selectSuggestedAction}
					stop={stop}
					attachments={attachments}
					setAttachments={setAttachments}
				/>
			</form>
		</div>
	)
}
