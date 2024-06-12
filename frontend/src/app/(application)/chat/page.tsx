import BreadCrumb from '~/components/breadcrumb'
import ChatCanvas from '~/components/chat/chat-canvas'
import ConversationsSidebar from '~/components/chat/conversation-list-sidebar'
import { Card } from '~/components/ui/card'
import { Heading } from '~/components/ui/heading'
import { Separator } from '~/components/ui/separator'

const ChatDashboard = () => {
	const breadcrumbItems = [{ title: 'Chat', link: '/chat' }]

	return (
		<div className="flex h-full flex-1 flex-col space-y-4 p-4 pt-6 md:p-8">
			<div className="flex flex-col gap-3">
				<BreadCrumb items={breadcrumbItems} />
				<div className="flex items-start justify-between">
					<Heading title={`Conversations`} description="" />
				</div>
				<Separator />
			</div>

			<div className="grid h-full grid-cols-7 gap-2">
				<Card className="col-span-2 h-full">
					<ConversationsSidebar />
				</Card>
				<Card className="col-span-5 h-full">
					<ChatCanvas />
				</Card>
			</div>
		</div>
	)
}

export default ChatDashboard
