import { ScrollArea } from '~/components/ui/scroll-area'
import Image from 'next/image'
import { Separator } from '../ui/separator'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '../ui/tabs'

const ConversationsSidebar = () => {
	const users = [
		{
			name: 'John Doe',
			status: 'Online',
			lastMessage: 'Hey, how are you?',
			unreadCount: 0,
			avatar: 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?q=80&w=3087&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D'
		},
		{
			name: 'Jane Smith',
			status: 'Offline',
			lastMessage: 'I will be there soon',
			unreadCount: 2,
			avatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?q=80&w=3087&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D'
		},
		{
			name: 'Alice Johnson',
			status: 'Online',
			lastMessage: 'What are you up to?',
			unreadCount: 1,
			avatar: 'https://images.unsplash.com/photo-1517841905240-472988babdf9?q=80&w=3087&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D'
		},
		{
			name: 'Bob Williams',
			status: 'Offline',
			lastMessage: 'See you tomorrow!',
			unreadCount: 0,
			avatar: 'https://images.unsplash.com/photo-1517841905240-472988babdf9?q=80&w=3087&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D'
		},
		{
			name: 'Emma Davis',
			status: 'Online',
			lastMessage: "Let's catch up soon",
			unreadCount: 3,
			avatar: 'https://images.unsplash.com/photo-1517841905240-472988babdf9?q=80&w=3087&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D'
		},
		{
			name: 'Michael Brown',
			status: 'Offline',
			lastMessage: 'Have a great day!',
			unreadCount: 0,
			avatar: 'https://images.unsplash.com/photo-1517841905240-472988babdf9?q=80&w=3087&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D'
		},
		{
			name: 'Olivia Wilson',
			status: 'Online',
			lastMessage: 'How was your weekend?',
			unreadCount: 2,
			avatar: 'https://images.unsplash.com/photo-1517841905240-472988babdf9?q=80&w=3087&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D'
		},
		{
			name: 'Sarah Johnson',
			status: 'Online',
			lastMessage: 'Good morning!',
			unreadCount: 1,
			avatar: 'https://images.unsplash.com/photo-1517841905240-472988babdf9?q=80&w=3087&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D'
		}
	].map(user => ({ ...user, avatar: user.avatar.replace(/w=\d+&h=\d+/, 'w=200&h=200') }))

	return (
		<ScrollArea
			className="flex h-full flex-col gap-2 px-2 py-4"
			key={'conversation_contacts_list_sidebar'}
		>
			<Tabs defaultValue="app-settings" className="w-full space-y-6">
				<TabsList className="flex w-full flex-row ">
					{['All', 'Unread', 'Unresolved'].map((tab, index) => {
						return (
							<TabsTrigger value={tab} className="flex-1" key={index}>
								{tab}
							</TabsTrigger>
						)
					})}
				</TabsList>
				<TabsContent value="All" className="space-y-4"></TabsContent>
			</Tabs>

			{users.map((user, index) => {
				return (
					<>
						<div
							key={index}
							className="my-auto flex flex-row items-center gap-4 px-3 py-2"
						>
							<Image
								src={user.avatar}
								height={50}
								width={50}
								className="object-fit aspect-square h-12 w-12 rounded-full"
								alt={`${user.name} avatar`}
							/>
							<div className="flex w-full flex-row justify-between">
								<div className="flex flex-col">
									<div className="flex items-center gap-2">
										<p className="text-sm"> {user.name}</p>
									</div>
									<p className="text-xs text-gray-500">{user.lastMessage}</p>
								</div>
								<div className="flex items-center justify-center">
									{user.unreadCount > 0 && (
										<div className="flex h-4 w-4 items-center justify-center rounded-full bg-primary text-xs text-white">
											{user.unreadCount}
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
