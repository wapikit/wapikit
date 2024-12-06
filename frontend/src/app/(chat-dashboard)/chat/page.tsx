import ChatCanvas from '~/components/chat/chat-canvas'
import ConversationsSidebar from '~/components/chat/conversation-list-sidebar'
import { Card } from '~/components/ui/card'

const ChatDashboard = () => {
	return (
		<div className="flex h-full flex-1 flex-col space-y-4 p-4 pt-6 md:p-8">
			<div className="grid h-screen grid-cols-7 gap-2">
				<Card className="col-span-2 h-full rounded-md">
					<ConversationsSidebar />
				</Card>
				<Card className="col-span-5 h-full rounded-md">
					<ChatCanvas />
				</Card>
			</div>
		</div>
	)
}

export default ChatDashboard
