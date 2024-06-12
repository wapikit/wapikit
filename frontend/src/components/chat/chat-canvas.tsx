import { ScrollArea } from '~/components/ui/scroll-area'
import { CardHeader, CardContent, CardFooter } from '../ui/card'
import { Separator } from '../ui/separator'
import { Input } from '../ui/input'
import { Button } from '../ui/button'
import { SendIcon, Image as ImageIcon } from 'lucide-react'
import Image from 'next/image'

const ChatCanvas = () => {
	const user = {
		name: 'John Doe',
		status: 'Online',
		lastMessage: 'Hey, how are you?',
		unreadCount: 0,
		avatar: 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?q=80&w=3087&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D'
	}

	return (
		<div className="relative flex h-full flex-col justify-between">
			<CardHeader className="flex flex-row p-3">
				<div className="flex flex-row gap-3">
					<Image
						src={user.avatar}
						height={50}
						width={50}
						className="object-fit aspect-square h-12 w-12 rounded-full"
						alt={`${user.name} avatar`}
					/>
					<div className="flex flex-col">
						<p className="text-base">{user.name}</p>
						{user.status === 'Online' ? (
							<span className="w-fit rounded-xl bg-primary-foreground px-1 text-center text-xs text-primary">
								Online
							</span>
						) : (
							<span>Last Active at {user.status}</span>
						)}
					</div>
				</div>
			</CardHeader>
			<Separator />

			<ScrollArea className="flex-1">
				<CardContent></CardContent>
			</ScrollArea>

			<CardFooter className="flex w-full flex-col gap-2">
				<Separator />
				<form className="flex w-full gap-2">
					<div className="flex items-center">
						<ImageIcon className="size-6" />
					</div>
					<Input placeholder="Type Message here" className="w-full" />
					<Button type="submit" className="rounded-full">
						<SendIcon className="size-4" />
					</Button>
				</form>
			</CardFooter>
		</div>
	)
}

export default ChatCanvas
