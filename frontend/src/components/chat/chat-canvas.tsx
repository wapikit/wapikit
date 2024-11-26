'use client'

import { ScrollArea } from '~/components/ui/scroll-area'
import { CardHeader, CardContent, CardFooter } from '../ui/card'
import { Separator } from '../ui/separator'
import { Input } from '../ui/input'
import { Button } from '../ui/button'
import { SendIcon, Image as ImageIcon, MoreVerticalIcon } from 'lucide-react'
import Image from 'next/image'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger
} from '../ui/dropdown-menu'
import { Icons } from '../icons'
import { useState } from 'react'
import {
	MessageDirectionEnum,
	MessageStatusEnum,
	MessageTypeEnum,
	useAssignConversation,
	useUnassignConversation
} from 'root/.generated'
import MessageRenderer from './message-renderer'
import { Sheet, SheetContent, SheetHeader, SheetTitle } from '../ui/sheet'
import { Label } from '../ui/label'
import { useRouter } from 'next/navigation'
import { errorNotification, successNotification } from '~/reusable-functions'

const ChatCanvas = () => {
	const [isBusy, setIsBusy] = useState(false)
	const [isContactInfoSheetVisible, setIsContactInfoSheetVisible] = useState(false)

	const router = useRouter()

	const assignConversationMutation = useAssignConversation()
	const unassignConversationMutation = useUnassignConversation()

	const user = {
		uniqueId: '12345',
		name: 'John Doe',
		status: 'Online',
		lastMessage: 'Hey, how are you?',
		unreadCount: 0,
		avatar: 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?q=80&w=3087&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D'
	}

	async function assignConversation() {
		try {
			if (isBusy) return

			setIsBusy(true)
			const assignConversationResponse = await assignConversationMutation.mutateAsync({
				data: {
					userId: ''
				},
				id: user.uniqueId
			})

			if (assignConversationResponse.data) {
				successNotification({
					message: 'Conversation assigned successfully'
				})
			} else {
				errorNotification({
					message: 'Failed to assign conversation'
				})
			}
		} catch (error) {
			console.error(error)
			errorNotification({
				message: 'Failed to assign conversation'
			})
		} finally {
			setIsBusy(false)
		}
	}

	async function unassignConversation() {
		try {
			if (isBusy) return

			setIsBusy(true)
			const unassignConversationResponse = await unassignConversationMutation.mutateAsync({
				data: {
					userId: ''
				},
				id: user.uniqueId
			})

			if (unassignConversationResponse.data) {
				successNotification({
					message: 'Conversation unassigned successfully'
				})
			} else {
				errorNotification({
					message: 'Failed to unassign conversation'
				})
			}
		} catch (error) {
			console.error(error)
			errorNotification({
				message: 'Failed to unassign conversation'
			})
		} finally {
			setIsBusy(false)
		}
	}

	const chatActions: {
		label: string
		icon: keyof typeof Icons
		onClick?: () => void
	}[] = [
		{
			label: 'Edit Contact',
			icon: 'edit',
			onClick: () => {
				router.push(`/contacts/new-or-edit/${user.uniqueId}`)
			}
		},
		{
			label: 'Assign to',
			icon: 'user'
		},
		{
			label: 'Unassign',
			icon: 'removeUser'
		},
		{
			label: 'Block',
			icon: 'xCircle'
		},
		{
			label: 'Mark As Resolved',
			icon: 'check'
		},
		{
			label: 'Info',
			icon: 'info',
			onClick: () => {
				setIsContactInfoSheetVisible(true)
			}
		}
	]

	return (
		<div className="relative flex h-full flex-col justify-between">
			<CardHeader className="item-center flex !flex-row justify-between rounded-t-md  bg-primary p-3 py-2">
				<div className="flex flex-row gap-3 ">
					<Image
						src={user.avatar}
						height={50}
						width={50}
						className="object-fit aspect-square h-10 w-10 rounded-full"
						alt={`${user.name} avatar`}
					/>
					<div className="flex flex-col">
						<p className="text-base text-primary-foreground">{user.name}</p>
						{user.status === 'Online' ? (
							<span className="w-fit rounded-xl bg-primary-foreground px-1 text-center text-xs text-primary">
								Online
							</span>
						) : (
							<span>Last Active at {user.status}</span>
						)}
					</div>
				</div>

				<div className="ml-auto">
					<DropdownMenu modal={false}>
						<DropdownMenuTrigger asChild>
							<MoreVerticalIcon className="text-bold h-5 w-5 text-primary-foreground" />
						</DropdownMenuTrigger>
						<DropdownMenuContent align="end">
							{chatActions.map((action, index) => {
								const Icon = Icons[action.icon]
								return (
									<DropdownMenuItem
										key={index}
										onClick={() => {
											if (action.onClick) {
												action.onClick()
											}
										}}
										className="flex flex-row items-center gap-2"
									>
										<Icon className="size-4" />
										{action.label}
									</DropdownMenuItem>
								)
							})}
						</DropdownMenuContent>
					</DropdownMenu>
				</div>
			</CardHeader>
			<Separator />

			<ScrollArea className="flex-1">
				<Sheet
					open={isContactInfoSheetVisible}
					onOpenChange={isOpen => {
						if (!isOpen) {
							setIsContactInfoSheetVisible(false)
						}
					}}
				>
					<SheetContent>
						<SheetHeader>
							<SheetTitle>Contact Info</SheetTitle>
						</SheetHeader>
						<div className="grid gap-4 py-4">
							{/* profile picture */}
							{/* user status */}
							{/* tags */}
							{/* list the user is in */}

							<div className="grid grid-cols-4 items-center gap-4">
								<Label htmlFor="name" className="text-right">
									Name
								</Label>
								<Input id="name" value="Pedro Duarte" className="col-span-3" />
							</div>
							<div className="grid grid-cols-4 items-center gap-4">
								<Label htmlFor="username" className="text-right">
									Username
								</Label>
								<Input id="username" value="@peduarte" className="col-span-3" />
							</div>
						</div>
						{/* <SheetFooter>
							<SheetClose asChild>
								<Button type="submit">Save changes</Button>
							</SheetClose>
						</SheetFooter> */}
					</SheetContent>
				</Sheet>

				<CardContent className="relative h-full w-full  bg-[#ebe5de] !py-4">
					<div className='absolute inset-0 z-20 h-full w-full  bg-[url("/assets/chat-canvas-bg.png")] bg-repeat opacity-20 ' />
					{Array(5)
						.fill(0)
						.map((_, index) => {
							// if odd then inbound else outbound

							const message = {
								content: 'Hii, hello world',
								conversationId: '1233453',
								createdAt: new Date().toISOString(),
								direction:
									index % 2 === 0
										? MessageDirectionEnum.InBound
										: MessageDirectionEnum.OutBound,
								message_type: MessageTypeEnum.Text,
								status: MessageStatusEnum.Read,
								uniqueId: `${12345 + index}`
							}

							return (
								<div className="relative z-30 w-full " key={index}>
									<MessageRenderer message={message} isActionsEnabled={true} />
								</div>
							)
						})}
				</CardContent>
			</ScrollArea>

			<CardFooter className="sticky bottom-0 z-30 flex w-full flex-col gap-2 bg-white">
				<Separator />
				<form className="flex w-full gap-2">
					<div className="flex items-center">
						<ImageIcon className="size-6" />
					</div>
					<Input placeholder="Type Message here" className="w-full" />
					<Button type="submit" className="rounded-full" disabled={isBusy}>
						<SendIcon className="size-4" />
					</Button>
				</form>
			</CardFooter>
		</div>
	)
}

export default ChatCanvas
