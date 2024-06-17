import { useEffect, useRef, useState } from 'react'
import { type z } from 'zod'
import { WebsocketEventDataMap, WebsocketEventEnum } from '../websocket-events'
import { generateUniqueId, getWebsocketUrl } from '~/reusable-functions'

export function useWebsocket(params: { token: string }) {
	const url = getWebsocketUrl(params.token)
	const [isConnected, setIsConnected] = useState(false)
	const [isConnecting] = useState(false)
	const wsRef = useRef<WebSocket | null>(null)
	const [pendingMessages] = useState<Map<string, (data: { status: 'ok' }) => void>>(new Map())

	useEffect(() => {
		wsRef.current = new WebSocket(url)
		wsRef.current.onopen = () => setIsConnected(true)
		wsRef.current.onclose = () => setIsConnected(false)
		wsRef.current.onmessage = event => {
			const message: z.infer<(typeof WebsocketEventDataMap)[WebsocketEventEnum]> = JSON.parse(
				event.data
			)

			const newParsedResponse = WebsocketEventDataMap[message.eventName].safeParse(message)

			if (newParsedResponse.success) {
				const parsedMessage = newParsedResponse.data

				console.log({ ...parsedMessage })
				switch (message.eventName) {
					case WebsocketEventEnum.MessageAcknowledgementEvent: {
						const resolve = pendingMessages.get(message.data.messageId)
						if (resolve) {
							resolve({ status: 'ok' })
							pendingMessages.delete(message.data.messageId)
						}

						break
					}

					case WebsocketEventEnum.MessageEvent: {
						// handle message event
						break
					}

					case WebsocketEventEnum.NotificationReadEvent: {
						// handle notification read event
						break
					}

					case WebsocketEventEnum.MessageReadEvent: {
						// handle message read event
						break
					}

					case WebsocketEventEnum.NewNotificationEvent: {
						// handle new notification event
						break
					}

					case WebsocketEventEnum.SystemReloadEvent: {
						// handle system reload event
						break
					}

					case WebsocketEventEnum.ConversationAssignmentEvent: {
						// handle conversation assignment event
						break
					}

					case WebsocketEventEnum.ConversationClosedEvent: {
						// handle conversation closed event
						break
					}

					case WebsocketEventEnum.NewConversationEvent: {
						// handle new conversation event
						break
					}

					default: {
						throw new Error('Unhandled event')
					}
				}
			} else {
				throw new Error('Invalid message')
			}
		}

		wsRef.current.onerror = error => {
			console.error(error)
			// reconnect try
		}
		return () => {
			wsRef.current?.close()
		}
	}, [pendingMessages, url])

	// async send function
	const sendMessage = async (
		payload: z.infer<(typeof WebsocketEventDataMap)[WebsocketEventEnum]>
	): Promise<{ status: 'ok' }> => {
		return new Promise(resolve => {
			const messageId = generateUniqueId()
			// add the event in pending messages
			pendingMessages.set(messageId, resolve)

			// send message here over websocket
			wsRef.current?.send(
				JSON.stringify({
					...payload
				})
			)
		})
	}

	return {
		isConnected,
		isConnecting,
		sendMessage
	}
}
