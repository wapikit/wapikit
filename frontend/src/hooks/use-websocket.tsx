import { useCallback, useEffect, useRef, useState } from 'react'
import { type z } from 'zod'
import { WebsocketEventDataMap, WebsocketEventEnum } from '../websocket-events'
import { generateUniqueId, getWebsocketUrl } from '~/reusable-functions'
import { useAuthState } from './use-auth-state'
import { WebsocketStatusEnum } from '~/types'
import { messageEventHandler } from '~/utils/websocket-handlers'
import { useConversationInboxStore } from '~/store/conversation-inbox.store'

const encoder = new TextEncoder()
const decoder = new TextDecoder('utf-8')

export function useWebsocket() {
	const [websocketStatus, setWebsocketStatus] = useState<WebsocketStatusEnum>(
		WebsocketStatusEnum.Idle
	)

	const { writeProperty, conversations } = useConversationInboxStore()

	const wsRef = useRef<WebSocket | null>(null)
	const [pendingMessages] = useState<Map<string, (data: { status: 'ok' }) => void>>(new Map())
	const { authState } = useAuthState()
	const sendWebsocketEvent = useCallback(
		async (
			payload: z.infer<(typeof WebsocketEventDataMap)[WebsocketEventEnum]>
		): Promise<{ status: 'ok' }> => {
			return new Promise(resolve => {
				const eventId = generateUniqueId()
				// add the event in pending messages
				pendingMessages.set(eventId, resolve)

				const eventInBinaryFormat = encoder.encode(
					JSON.stringify({
						...payload,
						eventId
					})
				)

				// send event here over websocket
				wsRef.current?.send(eventInBinaryFormat)
			})
		},
		[pendingMessages]
	)

	const _sendAcknowledgement = async (eventId: string) => {
		const data: z.infer<(typeof WebsocketEventDataMap)['MessageAcknowledgementEvent']> = {
			eventName: WebsocketEventEnum.MessageAcknowledgementEvent,
			eventId: eventId,
			data: {
				message: 'Acknowledged'
			}
		}
		await sendWebsocketEvent(data)
	}

	const tryConnectingToWebsocket = useCallback(() => {
		if (!authState.isAuthenticated || websocketStatus !== WebsocketStatusEnum.Idle) return
		setWebsocketStatus(() => WebsocketStatusEnum.Connecting)
		wsRef.current = new WebSocket(getWebsocketUrl(authState.data.token))
		wsRef.current.onopen = () => {
			setWebsocketStatus(() => WebsocketStatusEnum.Connected)
			setInterval(() => {
				const data: z.infer<(typeof WebsocketEventDataMap)['PingEvent']> = {
					eventName: WebsocketEventEnum.PingEvent,
					eventId: generateUniqueId(),
					data: {
						message: 'Ping!!!'
					}
				}
				sendWebsocketEvent(data).catch(error => console.error(error))
			}, 2000)
		}
		wsRef.current.onclose = () => setWebsocketStatus(() => WebsocketStatusEnum.Disconnected)

		wsRef.current.onmessage = async event => {
			const binaryData = event.data
			const buffer = binaryData instanceof Blob ? await binaryData.arrayBuffer() : binaryData
			const jsonString = decoder.decode(buffer)

			console.log({ jsonString })

			const message: z.infer<(typeof WebsocketEventDataMap)[WebsocketEventEnum]> =
				JSON.parse(jsonString)

			const schema = WebsocketEventDataMap[message.eventName]

			if (!schema) {
				console.log('unknown event received')
			}

			const newParsedResponse = schema.safeParse(message)

			let sendAcknowledgement = false

			if (newParsedResponse.success) {
				const parsedMessage = newParsedResponse.data
				switch (parsedMessage.eventName) {
					case WebsocketEventEnum.MessageAcknowledgementEvent: {
						const resolve = pendingMessages.get(parsedMessage.eventId)
						if (resolve) {
							resolve({ status: 'ok' })
							pendingMessages.delete(parsedMessage.eventId)
						}
						break
					}

					case WebsocketEventEnum.MessageEvent: {
						console.log('new message event received')
						const conversation = conversations.find(
							({ uniqueId }) => uniqueId === parsedMessage.data.conversationId
						)

						if (!conversation) {
							// ! TODO: this conversation is not in the store, fetch it from the server
							return
						}

						const done = await messageEventHandler({
							conversation: conversation,
							message: parsedMessage.data,
							writeProperty
						})
						sendAcknowledgement = done
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

				if (sendAcknowledgement) {
					await _sendAcknowledgement(parsedMessage.eventId)
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
	}, [authState, pendingMessages, sendWebsocketEvent])

	useEffect(() => {
		tryConnectingToWebsocket()

		return () => {
			console.log('closing websocket')
			wsRef.current?.close()
			setWebsocketStatus(WebsocketStatusEnum.Idle)
		}
	}, [tryConnectingToWebsocket])

	return {
		websocketStatus,
		sendMessage: sendWebsocketEvent
	}
}
