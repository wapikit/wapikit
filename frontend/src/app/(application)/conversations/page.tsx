'use client'

import { useEffect } from 'react'
import { useGetConversations } from 'root/.generated'
import { useConversationInboxStore } from '~/store/conversation-inbox.store'
import ChatCanvas from '~/components/chat/chat-canvas'
import ConversationsSidebar from '~/components/chat/conversation-list-sidebar'
import { Card } from '~/components/ui/card'
import { useSearchParams } from 'next/navigation'

const ChatDashboard = () => {
	const queryParams = useSearchParams()
	const conversationId = queryParams.get('id')

	const { writeProperty: writeConversationStoreProperty } = useConversationInboxStore()

	const { data: conversations } = useGetConversations({
		page: 1,
		per_page: 50
	})

	useEffect(() => {
		if (!conversations?.conversations) return

		writeConversationStoreProperty({
			conversations: conversations?.conversations || []
		})
	}, [conversations?.conversations, writeConversationStoreProperty])

	return (
		<div className="flex h-full flex-1 flex-col">
			<div className="grid h-screen grid-cols-7 gap-2 px-4">
				<Card className="col-span-2 h-full rounded-md">
					<ConversationsSidebar />
				</Card>
				<Card className="col-span-5 h-full rounded-md">
					<ChatCanvas conversationId={conversationId || undefined} />
				</Card>
			</div>
		</div>
	)
}

export default ChatDashboard
