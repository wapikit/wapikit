import { useCallback, useEffect, useRef, useState } from 'react'
import { type z } from 'zod'
import { WebsocketEventDataMap, WebsocketEventEnum } from '../websocket-events'
import { generateUniqueId, getWebsocketUrl } from '~/reusable-functions'
import { useAuthState } from './use-auth-state'
import { WebsocketStatusEnum } from '~/types'

export function useWebsocket() {
	const [websocketStatus, setWebsocketStatus] = useState<WebsocketStatusEnum>(
		WebsocketStatusEnum.Idle
	)
	const wsRef = useRef<WebSocket | null>(null)
	const [pendingMessages] = useState<Map<string, (data: { status: 'ok' }) => void>>(new Map())
	const { authState } = useAuthState()
	const sendMessage = useCallback(
		async (
			payload: z.infer<(typeof WebsocketEventDataMap)[WebsocketEventEnum]>
		): Promise<{ status: 'ok' }> => {
			return new Promise(resolve => {
				const messageId = generateUniqueId()
				// add the event in pending messages
				pendingMessages.set(messageId, resolve)

				// send message here over websocket
				wsRef.current?.send(
					JSON.stringify({
						...payload,
						messageId
					})
				)
			})
		},
		[pendingMessages]
	)

	const tryConnectingToWebsocket = useCallback(() => {
		if (!authState.isAuthenticated || websocketStatus !== WebsocketStatusEnum.Idle) return
		setWebsocketStatus(() => WebsocketStatusEnum.Connecting)
		wsRef.current = new WebSocket(getWebsocketUrl(authState.data.token))
		wsRef.current.onopen = () => {
			setWebsocketStatus(() => WebsocketStatusEnum.Connected)
			setInterval(() => {
				const data: z.infer<(typeof WebsocketEventDataMap)['PingEvent']> = {
					eventName: WebsocketEventEnum.PingEvent,
					data: {
						message: 'Ping!!!'
					}
				}
				sendMessage(data).catch(error => console.error(error))
			}, 2000)
		}
		wsRef.current.onclose = () => setWebsocketStatus(() => WebsocketStatusEnum.Disconnected)

		wsRef.current.onmessage = event => {
			const message: z.infer<(typeof WebsocketEventDataMap)[WebsocketEventEnum]> = JSON.parse(
				event.data
			)

			const schema = WebsocketEventDataMap[message.eventName]

			if (!schema) {
				console.log('unknown event received')
			}

			const newParsedResponse = schema.safeParse(message)

			// console.log({ erros: newParsedResponse.error?.errors, message })

			if (newParsedResponse.success) {
				// ! TODO: use this parsed message for type safety
				// const parsedMessage = newParsedResponse.data
				switch (message.eventName) {
					case WebsocketEventEnum.MessageAcknowledgementEvent: {
						const resolve = pendingMessages.get(message.messageId)
						if (resolve) {
							resolve({ status: 'ok' })
							pendingMessages.delete(message.messageId)
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

		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [authState, pendingMessages, sendMessage])

	useEffect(() => {
		// tryConnectingToWebsocket()

		return () => {
			console.log('closing websocket')
			wsRef.current?.close()
			setWebsocketStatus(WebsocketStatusEnum.Idle)
		}
	}, [tryConnectingToWebsocket])

	return {
		websocketStatus,
		sendMessage
	}
}
