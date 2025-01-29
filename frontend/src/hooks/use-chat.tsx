import { useCallback, useRef, useState } from 'react'
import { getBackendUrl } from '~/constants'
import { useAiChatStore } from '~/store/ai-chat-store'
import { useAuthState } from './use-auth-state'
import { ChatBotStateEnum } from '~/types'
import { AiChatMessageRoleEnum } from 'root/.generated'
import { errorNotification } from '~/reusable-functions'

const useChat = ({ chatId }: { chatId: string }) => {
	const { chats, updateChatMessage, pushMessage, currentChatMessages, updateUserMessageId } =
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

			let done = false
			while (!done) {
				const { done: readerDone, value } = await reader.read()
				done = readerDone
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
						// ! TODO: AI: error handling here
						setChatBotState(ChatBotStateEnum.Idle)
						currentMessageIdInStream.current = null
						return
					} else if (parsedChunk.type === 'messageDetails') {
						const userMessage = parsedChunk.userMessage
						const aiMessage = parsedChunk.aiMessage
						updateUserMessageId(userMessage.uniqueId)
						pushMessage({
							content: aiMessage.content,
							uniqueId: aiMessage.uniqueId,
							createdAt: aiMessage.createdAt,
							role: aiMessage.role
						})

						currentMessageIdInStream.current = parsedChunk.aiMessage.uniqueId
					}
				}

				// Keep the last incomplete chunk in the buffer
				buffer = chunks[chunks.length - 1]
			}
		},
		[pushMessage, updateChatMessage, updateUserMessageId]
	)

	const _sendAiMessage = useCallback(
		async (message: string) => {
			try {
				if (!authState.isAuthenticated) return
				// push the message before hand in the UI array, after that update the UI
				pushMessage({
					content: message,
					uniqueId: Math.random().toString(),
					createdAt: new Date().toISOString(),
					role: AiChatMessageRoleEnum.User
				})

				setInput('')

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

				if (response.status == 200) {
					await handleDataStream(reader)
				} else {
					errorNotification({
						message: 'Error during chat submission'
					})
				}
			} catch (error) {
				console.error('Error during chat submission:', error)
				setChatBotState(ChatBotStateEnum.Idle)
			}
		},
		[authState, currentChat?.uniqueId, handleDataStream, pushMessage]
	)

	const handleSubmit = useCallback(async () => {
		await _sendAiMessage(input)
	}, [_sendAiMessage, input])

	const selectSuggestedAction = (action: string) => {
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
		selectSuggestedAction
	}
}

export default useChat
