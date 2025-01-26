'use client'

import { useEffect, useRef } from 'react'
import { useGetAiChats } from 'root/.generated'
import { useAiChatStore } from '~/store/ai-chat-store'
import { useLayoutStore } from '~/store/layout.store'
import AiChatBox from '../ai/ai-chat-box'
import { useAuthState } from '~/hooks/use-auth-state'

const AiChatProvider = ({ children }: { children: React.ReactNode }) => {
	const { featureFlags } = useLayoutStore()
	const { writeProperty } = useAiChatStore()

	const { authState } = useAuthState()

	const isSetupDone = useRef(false)

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
		if (isSetupDone.current) {
			return
		}

		writeProperty({
			chats: chats?.chats || []
		})

		isSetupDone.current = true
	}, [writeProperty, chats?.chats])

	return (
		<>
			<AiChatBox />
			{children}
		</>
	)
}

export default AiChatProvider
