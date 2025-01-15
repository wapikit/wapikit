import { useCallback, useRef, useState } from 'react'
import { getBackendUrl } from '~/constants'
import { useAiChatStore } from '~/store/ai-chat-store'
import { useAuthState } from './use-auth-state'
import { ChatBotStateEnum } from '~/types'

const useChat = ({ chatId }: { chatId: string }) => {
	const { chats, updateChatMessage, pushMessage, currentChatMessages, editMessage } =
		useAiChatStore()
	const { authState } = useAuthState()
	const [input, setInput] = useState('')
	const currentMessageIdInStream = useRef<string | null>(null)
	const [chatBotState, setChatBotState] = useState<ChatBotStateEnum>(ChatBotStateEnum.Idle)
	const currentChat = chats.find(chat => chat.uniqueId === chatId)

	const handleDataStream = useCallback(
		async (reader: ReadableStreamDefaultReader) => {
			const decoder = new TextDecoder()
			let buffer = ''

			while (true) {
				const { done, value } = await reader.read()
				if (done) break

				buffer += decoder.decode(value, { stream: true })
				const chunks = buffer.split('\n')

				// Process each chunk
				for (let i = 0; i < chunks.length - 1; i++) {
					const chunk = chunks[i]
					const parsedChunk = JSON.parse(chunk)

					if (parsedChunk.type === 'text-delta') {
						if (!currentMessageIdInStream.current) return
						updateChatMessage(currentMessageIdInStream.current, parsedChunk.content)
					} else if (parsedChunk.type === 'finish') {
						setChatBotState(ChatBotStateEnum.Idle)
						return
					} else if (parsedChunk.type === 'messageDetails') {
						currentMessageIdInStream.current = parsedChunk.uniqueId
						// push the message to the current message list
						pushMessage({
							uniqueId: parsedChunk.uniqueId,
							content: parsedChunk.content,
							createdAt: parsedChunk.createdAt,
							role: parsedChunk.role
						})
					}
				}

				// Keep the last incomplete chunk in the buffer
				buffer = chunks[chunks.length - 1]
			}
		},
		[chatId, updateChatMessage]
	)

	const _sendAiMessage = useCallback(
		async (message: string) => {
			try {
				if (!authState.isAuthenticated) return
				setChatBotState(ChatBotStateEnum.Streaming)
				const response = await fetch(
					`${getBackendUrl()}/ai/chat/${currentChat?.uniqueId}/messages`,
					{
						body: JSON.stringify({ query: message }),
						method: 'POST',
						headers: {
							Accept: 'application/json',
							'x-access-token': authState.data.token || '',
							'Content-Type': 'application/json'
						},
						cache: 'no-cache',
						credentials: 'include'
					}
				)

				if (!response.body) throw new Error('No response body')
				const reader = response.body.getReader()
				await handleDataStream(reader)
			} catch (error) {
				console.error('Error during chat submission:', error)
				setChatBotState(ChatBotStateEnum.Idle)
			}
		},
		[authState, currentChat, handleDataStream]
	)

	const handleSubmit = useCallback(async () => {
		await _sendAiMessage(input)
	}, [_sendAiMessage, input])

	function selectSuggestedAction(action: string) {
		_sendAiMessage(action).catch(error => console.error(error))
	}

	return {
		currentChat,
		chatBotState,
		handleSubmit,
		currentMessageIdInStream,
		currentChatMessages,
		setInput,
		input,
		selectSuggestedAction,
	}
}

export default useChat
