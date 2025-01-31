import { useCallback, useEffect, useState } from 'react'
import { getBackendUrl } from '~/constants'
import { useAuthState } from '~/hooks/use-auth-state' // Assuming you have an auth hook
import { errorNotification } from '~/reusable-functions'

const useServiceSideEvents = () => {
	const { authState } = useAuthState()
	const [connectionState, setConnectionState] = useState<
		'Connecting' | 'Connected' | 'Disconnected'
	>('Disconnected')

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

				console.log({ parsedChunk })
			}

			// Keep the last incomplete chunk in the buffer
			buffer = chunks[chunks.length - 1]
		}
	}, [])

	const connectToSseEndpoint = useCallback(async () => {
		if (!authState?.isAuthenticated || !authState?.data?.token) {
			console.warn('User is not authenticated. SSE will not connect.')
			return
		}

		if (connectionState === 'Connected' || connectionState === 'Connecting') {
			console.warn('Already connected or connecting. Skipping...')
			return
		}

		const response = await fetch(`${getBackendUrl()}/events`, {
			method: 'GET',
			headers: {
				Accept: 'text/event-stream',
				'x-access-token': authState.data.token || '',
				'Content-Type': 'application/json'
			},
			cache: 'no-cache',
			credentials: 'include'
		})

		if (!response.body) throw new Error('No response body')
		const reader = response.body.getReader()

		if (response.status == 200) {
			setConnectionState(() => 'Connected')
			await handleDataStream(reader)
		} else {
			errorNotification({
				message: 'Error during listening to SSE. Please contact support immediately.'
			})
		}
	}, [authState, connectionState, handleDataStream])

	useEffect(() => {
		connectToSseEndpoint().catch(error => console.error('Error during SSE connection:', error))

		return () => {
			// Cleanup
		}
	}, [authState, authState, connectToSseEndpoint])

	return { connectionState }
}

export default useServiceSideEvents
