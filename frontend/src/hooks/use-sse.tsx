import { useEffect, useRef, useState } from 'react'
import { getBackendUrl } from '~/constants'
import { useAuthState } from '~/hooks/use-auth-state'

const MAX_RECONNECT_ATTEMPTS = 5
const RECONNECT_INTERVAL = 5000 // 5 seconds

const useServerSideEvents = () => {
	const { authState } = useAuthState()
	const [connectionState, setConnectionState] = useState<
		'Connecting' | 'Connected' | 'Disconnected'
	>('Disconnected')

	const eventSourceRef = useRef<EventSource | null>(null)
	const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null)
	const reconnectAttemptsRef = useRef(0)

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
			setConnectionState('Connecting')

			const sseUrl = `${getBackendUrl()}/events?token=${authState.data.token}`
			eventSourceRef.current = new EventSource(sseUrl, {
				withCredentials: true
			})

			eventSourceRef.current.onopen = () => {
				console.log('SSE connection established')
				reconnectAttemptsRef.current = 0
				setConnectionState('Connected')
			}

			eventSourceRef.current.onmessage = event => {
				try {
					const parsedData = JSON.parse(event.data)
					console.log('SSE message received:', parsedData)
				} catch (error) {
					console.error('SSE message parsing error:', error)
				}
			}

			eventSourceRef.current.onerror = error => {
				console.error('SSE connection error:', error)
				setConnectionState('Disconnected')

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
			setConnectionState('Disconnected')
		}
	}, [authState]) // Reconnect only when auth state changes

	return { connectionState }
}

export default useServerSideEvents
