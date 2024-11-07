'use client'

import { useWebsocket } from '~/hooks/use-websocket'

const WebsocketConnectionProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
	useWebsocket()
	return children
}

export default WebsocketConnectionProvider
