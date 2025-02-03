import { useCallback, useRef, useState } from 'react'
import { getBackendUrl } from '~/constants'
import { useAuthState } from './use-auth-state'

const useAiQuery = (
	params:
		| { purpose: 'segment-recommendation'; contactId?: string; conversationId?: string }
		| {
				purpose: 'conversation-summary'
				conversationId: string
		  }
) => {
	const { authState } = useAuthState()
	const currentMessageIdInStream = useRef<string | null>(null)
	const [responseState, setResponseState] = useState<'streaming' | 'idle' | 'errored'>('idle')

	const [completeResponse, setCompleteResponse] = useState<string>('')

	const handleDataStream = useCallback(async (reader: ReadableStreamDefaultReader) => {
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
					setCompleteResponse(data => data + parsedChunk.content)
				} else if (parsedChunk.type === 'finish') {
					setResponseState('idle')
					currentMessageIdInStream.current = null
					return
				} else {
					// * IMPOSSIBLE CASE
				}
			}

			// Keep the last incomplete chunk in the buffer
			buffer = chunks[chunks.length - 1]
		}
	}, [])

	const sendQuery = useCallback(async () => {
		try {
			if (!authState.isAuthenticated) return
			setResponseState('streaming')
			const backendEndpoint =
				params.purpose === 'conversation-summary'
					? 'chat-summary'
					: 'segment-recommendation'
			let body = JSON.stringify({ chatId: params.conversationId })
			if (params.purpose === 'segment-recommendation') {
				body = JSON.stringify({
					...(params.contactId && { contactId: params.contactId }),
					...(params.conversationId && { chatId: params.conversationId })
				})
			} else {
				// * defined above already
			}

			const response = await fetch(`${getBackendUrl()}/ai/${backendEndpoint}`, {
				body,
				method: 'POST',
				headers: {
					Accept: 'application/json',
					'x-access-token': authState.data.token || '',
					'Content-Type': 'application/json'
				},
				cache: 'no-cache',
				credentials: 'include'
			})

			if (!response.body) throw new Error('No response body')
			const reader = response.body.getReader()

			if (params.purpose === 'conversation-summary') {
				await handleDataStream(reader)
			} else {
				// * this is a non steaming response
				const json = await response.json()
				setCompleteResponse(json)
			}
		} catch (error) {
			console.error('Error during chat submission:', error)
			setResponseState('idle')
		}
	}, [authState, handleDataStream, params])

	return {
		responseState,
		sendQuery,
		currentMessageIdInStream,
		completeResponse
	}
}

export default useAiQuery
