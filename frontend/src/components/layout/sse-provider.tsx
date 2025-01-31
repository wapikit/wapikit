'use client'

import useServerSideEvents from '~/hooks/use-sse'

const SseConnectionProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
	useServerSideEvents()

	return children
}

export default SseConnectionProvider
