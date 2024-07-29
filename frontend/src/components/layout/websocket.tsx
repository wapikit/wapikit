'use client'

import { useWebsocket } from '~/hooks/use-websocket'

const WebsocketConnectionProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
	const { websocketStatus } = useWebsocket()
	console.log({ websocketStatus })

	return children
}

export default WebsocketConnectionProvider
