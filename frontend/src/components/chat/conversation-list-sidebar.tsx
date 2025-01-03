import { ScrollArea } from '~/components/ui/scroll-area'
import Image from 'next/image'
import { Separator } from '../ui/separator'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '../ui/tabs'
import { useConversationInboxStore } from '~/store/conversation-inbox.store'
import { Icons } from '../icons'

const ConversationsSidebar = () => {
	const { conversations, writeProperty: writeConversationStoreProperty } =
		useConversationInboxStore()

	return (
		<ScrollArea
			className="flex h-full flex-col gap-2 px-4 py-4"
			key={'conversation_contacts_list_sidebar'}
		>
			<Tabs defaultValue="All" className="w-full space-y-6">
				<TabsList className="flex w-full flex-row ">
					{['All', 'Unread', 'Unresolved'].map((tab, index) => {
						return (
							<TabsTrigger
								value={tab}
								className="flex flex-1 items-center gap-1"
								key={index}
							>
								{
									{
										All: <Icons.message className="size-4" />,
										Unread: <Icons.bell className="size-4" />,
										Unresolved: <Icons.help className="size-4" />
									}[tab]
								}
								{tab}
							</TabsTrigger>
						)
					})}
				</TabsList>
				<TabsContent value="All" className="space-y-4"></TabsContent>
			</Tabs>

			{conversations.length === 0 && (
				<div className="flex h-full flex-col items-center justify-center">
					<Icons.message className="size-6 font-normal text-muted-foreground" />
					<p className="text-gray-500">No conversations yet</p>
				</div>
			)}

			{conversations.map((conversation, index) => {
				const lastMessage =
					typeof conversation.messages.at(-1)?.content === 'string' &&
					conversation.messages.at(-1)?.content
						? conversation.messages.at(-1)?.content
						: ''

				return (
					<>
						<div
							key={index}
							className="my-auto flex cursor-pointer flex-row items-center gap-4 px-3 py-2 hover:bg-gray-100"
							onClick={() => {
								writeConversationStoreProperty({
									currentConversation: conversation
								})
							}}
						>
							<Image
								src={'/assets/empty-pfp.png'}
								height={50}
								width={50}
								className="object-fit aspect-square h-12 w-12 rounded-full"
								alt={`${conversation.contact.uniqueId}-avatar`}
							/>

							<div className="flex w-full flex-row justify-between">
								<div className="flex flex-col">
									<div className="flex items-center gap-2">
										<p className="text-sm"> {conversation.contact.name}</p>
									</div>
									<p className="text-xs text-gray-500">{lastMessage || ''}</p>
								</div>
								<div className="flex items-center justify-center">
									{conversation.numberOfUnreadMessages > 0 && (
										<div className="flex h-4 w-4 items-center justify-center rounded-full bg-primary text-xs text-white">
											{conversation.numberOfUnreadMessages}
										</div>
									)}
								</div>
							</div>
						</div>
						<Separator className="mx-1" />
					</>
				)
			})}
		</ScrollArea>
	)
}

export default ConversationsSidebar
