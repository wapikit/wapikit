'use client'

import type { Attachment } from 'ai'
import { useState } from 'react'
import { ChatHeader } from '~/components/ai/chat-header'
import { AiChatInput } from './ai-chat-input'
import { Messages } from './messages'
import { AiChatSchema, useGetAiChatMessageVotes } from 'root/.generated'
import useChat from '~/hooks/use-chat'
import { ChatBotStateEnum } from '~/types'

export function AiChat({ chat }: { chat: AiChatSchema }) {
	const { handleSubmit, chatBotState, input, setInput, selectSuggestedAction } = useChat({
		chatId: chat.uniqueId
	})

	// ! TODO: handle pagination here
	const { data: votes } = useGetAiChatMessageVotes(chat.uniqueId, {
		page: 1,
		per_page: 100
	})

	const [attachments, setAttachments] = useState<Array<Attachment>>([])

	return (
		<>
			<div className="flex h-dvh w-full min-w-0 flex-col bg-background">
				<ChatHeader chatTitle={chat.title} />
				<Messages
					chatId={chat.uniqueId}
					isLoading={chatBotState === ChatBotStateEnum.Streaming}
					votes={votes?.votes}
					isReadonly={false}
				/>

				<form className="mx-auto flex w-full gap-2 bg-background pb-4 md:pb-6">
					<AiChatInput
						chatId={chat.uniqueId}
						input={input}
						setInput={setInput}
						handleSubmit={handleSubmit}
						isLoading={chatBotState === ChatBotStateEnum.Streaming}
						selectSuggestedAction={selectSuggestedAction}
						stop={stop}
						attachments={attachments}
						setAttachments={setAttachments}
					/>
				</form>
			</div>
		</>
	)
}
