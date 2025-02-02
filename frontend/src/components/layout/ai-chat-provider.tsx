'use client'

import { useEffect } from 'react'
import { useGetAiChats } from 'root/.generated'
import { useAiChatStore } from '~/store/ai-chat-store'
import { useLayoutStore } from '~/store/layout.store'
import { useAuthState } from '~/hooks/use-auth-state'

const AiChatProvider = ({ children }: { children: React.ReactNode }) => {
	const { featureFlags } = useLayoutStore()
	const { writeProperty } = useAiChatStore()

	const { authState } = useAuthState()

	// ! TODO: handle the pagination here
	const { data: chats } = useGetAiChats(
		{
			page: 1,
			per_page: 20
		},
		{
			query: {
				enabled:
					!!authState.isAuthenticated &&
					featureFlags?.SystemFeatureFlags.isAiIntegrationEnabled
			}
		}
	)

	useEffect(() => {
		writeProperty({
			chats: chats?.chats || []
		})
	}, [writeProperty, chats?.chats])

	return <>{children}</>
}

export default AiChatProvider
