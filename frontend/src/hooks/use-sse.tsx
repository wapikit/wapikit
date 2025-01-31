import { useEffect, useRef, useState } from 'react'
import { getBackendUrl } from '~/constants'
import { useAuthState } from '~/hooks/use-auth-state'
import { useConversationInboxStore } from '~/store/conversation-inbox.store'
import { messageEventHandler } from '~/utils/sse-handlers'
import { ApiServerEventDataMap, ApiServerEventEnum } from '~/api-server-events'
import { SseEventSourceStateEnum } from '~/types'

const MAX_RECONNECT_ATTEMPTS = 5
const RECONNECT_INTERVAL = 5000 // 5 seconds

const useServerSideEvents = () => {
	const { authState } = useAuthState()
	const [connectionState, setConnectionState] = useState<SseEventSourceStateEnum>(
		SseEventSourceStateEnum.Disconnected
	)

	const eventSourceRef = useRef<EventSource | null>(null)
	const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null)
	const reconnectAttemptsRef = useRef(0)

	const { writeProperty, conversations } = useConversationInboxStore()
	const conversationsRef = useRef(conversations)
	useEffect(() => {
		conversationsRef.current = conversations
	}, [conversations])

	useEffect(() => {
		// Skip if already connected or connecting
		if (eventSourceRef.current) return

		// Only connect when properly authenticated
		if (!authState?.isAuthenticated || !authState?.data?.token) {
			return
		}

		const connectToSseEndpoint = () => {
			console.log(
				`Attempting SSE connection (attempt ${reconnectAttemptsRef.current + 1}/${MAX_RECONNECT_ATTEMPTS})`
			)
			setConnectionState(SseEventSourceStateEnum.Connecting)

			const sseUrl = `${getBackendUrl()}/events?token=${authState.data.token}`
			eventSourceRef.current = new EventSource(sseUrl, {
				withCredentials: true
			})

			eventSourceRef.current.onopen = () => {
				console.log('SSE connection established')
				reconnectAttemptsRef.current = 0
				setConnectionState(SseEventSourceStateEnum.Connected)
			}

			eventSourceRef.current.addEventListener('NewMessage', event => {
				console.log('SSE message received:', event)
				const schema = ApiServerEventDataMap[ApiServerEventEnum.NewMessageEvent]

				console.log('Message data:', event.data)

				const parsedMessageData = schema.safeParse(JSON.parse(event.data))

				if (parsedMessageData.success === false) {
					console.error('Failed to parse message data:', parsedMessageData.error)
					return
				}

				console.log({
					parsedMessageData: parsedMessageData.data
				})

				messageEventHandler({
					conversations: conversationsRef.current,
					eventData: parsedMessageData.data,
					writeProperty
				})
			})

			eventSourceRef.current.onerror = error => {
				console.error('SSE connection error:', error)
				setConnectionState(SseEventSourceStateEnum.Disconnected)

				// Cleanup current connection
				eventSourceRef.current?.close()
				eventSourceRef.current = null

				// Manage reconnect attempts
				reconnectAttemptsRef.current++

				if (reconnectAttemptsRef.current >= MAX_RECONNECT_ATTEMPTS) {
					console.error('Maximum reconnect attempts reached. Stopping SSE.')
					return
				}

				// Schedule new reconnect attempt
				if (reconnectTimeoutRef.current) clearTimeout(reconnectTimeoutRef.current)
				reconnectTimeoutRef.current = setTimeout(connectToSseEndpoint, RECONNECT_INTERVAL)
			}
		}

		connectToSseEndpoint()

		return () => {
			// Cleanup on unmount or auth state change
			if (reconnectTimeoutRef.current) {
				clearTimeout(reconnectTimeoutRef.current)
				reconnectTimeoutRef.current = null
			}
			if (eventSourceRef.current) {
				eventSourceRef.current.close()
				eventSourceRef.current = null
			}
			reconnectAttemptsRef.current = 0
			setConnectionState(SseEventSourceStateEnum.Disconnected)
		}
	}, [authState, writeProperty]) // Reconnect only when auth state changes

	return { connectionState }
}

export default useServerSideEvents
