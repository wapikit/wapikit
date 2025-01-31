'use client'

import useServerSideEvents from '~/hooks/use-sse'

const WebsocketConnectionProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
	// useWebsocket()
	useServerSideEvents('/api/events')
	return children
}

export default WebsocketConnectionProvider
