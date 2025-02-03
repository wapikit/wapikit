'use client'

import { ScrollArea } from '~/components/ui/scroll-area'
import { CardHeader, CardFooter } from '../ui/card'
import { Separator } from '../ui/separator'
import { Input } from '../ui/input'
import { Button } from '~/components/ui/button'
import { SendIcon, Image as ImageIcon, MoreVerticalIcon } from 'lucide-react'
import { motion } from 'framer-motion'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger
} from '../ui/dropdown-menu'
import { Icons } from '../icons'
import { useCallback, useEffect, useRef, useState } from 'react'
import {
	type ConversationSchema,
	MessageTypeEnum,
	useAssignConversation,
	useGetConversationResponseSuggestions,
	useGetOrganizationMembers,
	useSendMessageInConversation,
	useUnassignConversation
} from 'root/.generated'
import MessageRenderer from './message-renderer'
import { useRouter } from 'next/navigation'
import { errorNotification, successNotification } from '~/reusable-functions'
import { Modal } from '../ui/modal'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '../ui/form'
import { Select, SelectContent, SelectItem, SelectTrigger } from '../ui/select'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { AssignConversationForm } from '~/schema'
import { type z } from 'zod'
import { isPresent } from 'ts-is-present'
import { useLayoutStore } from '~/store/layout.store'
import ContactDetailsSheet from '../contact-details-sheet'
import { useConversationInboxStore } from '~/store/conversation-inbox.store'
import Image from 'next/image'
import { useScrollToBottom } from '~/hooks/use-scroll-to-bottom'
import { SparklesIcon } from '../ai/icons'

const ChatCanvas = ({ conversationId }: { conversationId?: string }) => {
	const [isBusy, setIsBusy] = useState(false)
	const [isConversationAssignModalOpen, setIsConversationAssignModalOpen] = useState(false)
	const { conversations } = useConversationInboxStore()

	const inputRef = useRef<HTMLInputElement>(null)

	const currentConversation = conversations.find(
		conversation => conversation.uniqueId === conversationId
	)

	const [messagesContainerRef, messagesEndRef] = useScrollToBottom<HTMLDivElement>()

	const router = useRouter()
	const { writeProperty } = useLayoutStore()
	const { writeProperty: writeConversationInboxStoreProperty } = useConversationInboxStore()

	const assignConversationMutation = useAssignConversation()
	const unassignConversationMutation = useUnassignConversation()
	const sendMessageInConversation = useSendMessageInConversation()

	const assignConversationForm = useForm<z.infer<typeof AssignConversationForm>>({
		resolver: zodResolver(AssignConversationForm)
	})

	const { data: organizationMembersResponse } = useGetOrganizationMembers({
		page: 1,
		per_page: 50,
		sortBy: 'asc'
	})

	async function assignConversation(data: z.infer<typeof AssignConversationForm>) {
		try {
			if (isBusy || !currentConversation) return

			setIsBusy(true)
			const assignConversationResponse = await assignConversationMutation.mutateAsync({
				data: {
					organizationMemberId: data.assignee
				},
				id: currentConversation.uniqueId
			})

			if (assignConversationResponse.data) {
				successNotification({
					message: 'Conversation assigned successfully'
				})
				setIsConversationAssignModalOpen(false)
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
			if (isBusy || !currentConversation) return

			setIsBusy(true)
			const unassignConversationResponse = await unassignConversationMutation.mutateAsync({
				data: {
					userId: ''
				},
				id: currentConversation.uniqueId
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
				router.push(`/contacts/new-or-edit/${currentConversation?.uniqueId}`)
			}
		},
		{
			label: 'Unassign',
			icon: 'removeUser',
			onClick() {
				unassignConversation().catch(error => console.error(error))
			}
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
				writeProperty({
					contactSheetContactId: currentConversation?.contact.uniqueId
				})
			}
		}
	]

	const [messageContent, setMessageContent] = useState<string | null>(null)

	const {
		data: suggestions,
		refetch: refetchSuggestions,
		isFetching: isFetchingSuggestions,
		isRefetching: isRefetchingSuggestions
	} = useGetConversationResponseSuggestions(
		{
			conversationId: currentConversation?.uniqueId || ''
		},
		{
			query: {
				enabled: !!currentConversation
			}
		}
	)

	const sendMessage = useCallback(
		async (message: string) => {
			try {
				console.log('sendMessage', {
					currentConversation,
					message
				})

				if (!currentConversation || !message) return

				setIsBusy(true)

				const sendMessageResponse = await sendMessageInConversation.mutateAsync({
					data: {
						messageData: {
							text: message
						},
						messageType: MessageTypeEnum.Text
					},
					id: currentConversation.uniqueId
				})

				if (sendMessageResponse.message) {
					const conversation = conversations.find(
						convo => convo.uniqueId === sendMessageResponse.message.conversationId
					)

					if (!conversation) {
						return false
					}

					const updatedConversation: ConversationSchema = {
						...conversation,
						messages: [...conversation.messages, sendMessageResponse.message]
					}

					writeConversationInboxStoreProperty({
						conversations: conversations.map(convo =>
							convo.uniqueId === conversation.uniqueId ? updatedConversation : convo
						)
					})

					console.log('Message sent successfully')

					setMessageContent(() => null)
				} else {
					errorNotification({
						message: 'Failed to send message'
					})
				}
			} catch (error) {
				console.error(error)
				errorNotification({
					message: 'Failed to send message'
				})
			} finally {
				setIsBusy(false)
			}
		},
		[
			currentConversation,
			sendMessageInConversation,
			conversations,
			writeConversationInboxStoreProperty
		]
	)

	useEffect(() => {
		// check if input is focussed, on enter sendMessage function should be called
		const handleKeyDown = (event: KeyboardEvent) => {
			if (
				document.activeElement === inputRef.current &&
				event.key === 'Enter' &&
				messageContent
			) {
				sendMessage(messageContent).catch(error => console.error(error))
			}
		}

		inputRef.current?.addEventListener('keydown', handleKeyDown)
	}, [messageContent, sendMessage])

	return (
		<div className="relative flex h-full flex-col justify-between">
			<ContactDetailsSheet />

			<Modal
				title="Assign Conversation to"
				description="Select a team member to assign this conversation to."
				isOpen={isConversationAssignModalOpen}
				onClose={() => {
					setIsConversationAssignModalOpen(false)
				}}
			>
				<div className="flex w-full items-center justify-end space-x-2 pt-6">
					<Form {...assignConversationForm}>
						<form
							onSubmit={assignConversationForm.handleSubmit(assignConversation)}
							className="w-full space-y-8"
						>
							<div className="flex flex-col gap-8">
								<FormField
									control={assignConversationForm.control}
									name="assignee"
									render={({ field }) => (
										<FormItem>
											<FormLabel className="flex flex-row items-center gap-2">
												Select Assignee
											</FormLabel>
											<FormControl>
												<Select
													disabled={isBusy}
													onValueChange={e => {
														field.onChange(e)
													}}
													name="assignee"
												>
													<SelectTrigger>
														<div>
															{organizationMembersResponse?.members
																?.map(member => {
																	if (
																		member.uniqueId ===
																		assignConversationForm.getValues(
																			'assignee'
																		)
																	) {
																		const stringToReturn = `${member.name} - ${member.email}`
																		return stringToReturn
																	} else {
																		return null
																	}
																})
																.filter(isPresent)[0] ||
																'Select Assignee'}
														</div>
													</SelectTrigger>
													<SelectContent
														side="bottom"
														className="max-h-64"
													>
														{!organizationMembersResponse ||
														organizationMembersResponse?.members
															.length === 0 ? (
															<SelectItem
																value={'no message template'}
																disabled
															>
																No organization member.
															</SelectItem>
														) : (
															<>
																{organizationMembersResponse?.members.map(
																	member => (
																		<SelectItem
																			key={`${member.uniqueId}`}
																			value={member.uniqueId}
																		>
																			{member.name} -{' '}
																			{member.email}
																		</SelectItem>
																	)
																)}
															</>
														)}
													</SelectContent>
												</Select>
											</FormControl>
											<FormMessage />
										</FormItem>
									)}
								/>
							</div>
							<Button disabled={isBusy} className="ml-auto mr-0 w-full" type="submit">
								Assign Conversation
							</Button>
						</form>
					</Form>
				</div>
			</Modal>

			{currentConversation ? (
				<>
					<CardHeader className="item-center flex !flex-row justify-between rounded-t-md  bg-primary p-3 py-2 dark:bg-[#202c33]">
						<div className="flex flex-row items-center gap-3">
							<Image
								src={'/assets/empty-pfp.png'}
								height={50}
								width={50}
								className="object-fit aspect-square h-10 w-10 cursor-pointer rounded-full"
								alt={`${currentConversation.uniqueId} avatar`}
								onClick={() => {
									writeProperty({
										contactSheetContactId: currentConversation?.contact.uniqueId
									})
								}}
							/>
							<p className="align-middle text-base text-primary-foreground">
								{currentConversation.contact.name}
							</p>
						</div>

						<div className="ml-auto flex flex-row items-center gap-4">
							<div className="flex flex-row items-center gap-2 text-sm text-secondary">
								Assigned To:
								{currentConversation.assignedTo ? (
									// ! TODO: show tippy on hover with user details
									<span>{currentConversation.assignedTo.name}</span>
								) : (
									<span
										className="rounded-full bg-white p-1"
										onClick={() => {
											setIsConversationAssignModalOpen(() => true)
										}}
									>
										<Icons.add className="size-4 text-secondary-foreground" />
									</span>
								)}
							</div>
							<DropdownMenu modal={false}>
								<DropdownMenuTrigger asChild>
									<MoreVerticalIcon className="text-bold h-5 w-5  text-secondary" />
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

					{/* ! TODO: this should always open at the end of scroll container */}
					<ScrollArea className="h-screen bg-[#ebe5de] !py-4 px-2 !pb-64 dark:bg-[#111b21]">
						<div className='absolute inset-0 z-20 h-full w-full  bg-[url("/assets/chat-canvas-bg.png")] bg-repeat opacity-20' />
						<div className="flex h-full flex-col gap-1">
							{currentConversation.messages.map((message, index) => {
								return (
									<div
										className="relative z-30 w-full"
										key={index}
										ref={messagesContainerRef}
									>
										<MessageRenderer
											message={message}
											isActionsEnabled={true}
										/>
										<div
											ref={messagesEndRef}
											className="min-h-[24px] min-w-[24px] shrink-0"
										/>
									</div>
								)
							})}
						</div>
					</ScrollArea>

					<div className="sticky bottom-0 z-30 ">
						{suggestions?.suggestions?.length ? (
							<motion.div
								initial={{ opacity: 0, y: 20 }}
								animate={{ opacity: 1, y: 0 }}
								exit={{ opacity: 0, y: 20 }}
								transition={{ duration: 0.5 }}
								className="relative z-30 w-full"
							>
								<div className="flex flex-1 flex-row justify-start gap-4 overflow-scroll px-5 pb-4">
									{suggestions?.suggestions.map((response, index) => {
										if (index > 3) return null
										return (
											<div
												key={index}
												className="h-auto w-fit cursor-pointer justify-start  whitespace-normal rounded-md bg-primary p-1 px-4 py-2 text-left text-sm font-medium text-primary-foreground "
												onClick={() => {
													setMessageContent(() => response)
													sendMessage(response).catch(error =>
														console.error(error)
													)
												}}
											>
												{response}
											</div>
										)
									})}
								</div>
							</motion.div>
						) : null}
						<CardFooter className="flex w-full flex-col gap-2 bg-white dark:bg-[#202c33]">
							<Separator />
							<form
								className="flex h-full w-full gap-2 "
								onSubmit={e => {
									e.preventDefault()
									sendMessage(messageContent || '').catch(error =>
										console.error(error)
									)
								}}
							>
								<div className="flex items-center">
									<ImageIcon className="size-6" />
								</div>
								<Input
									placeholder="Type Message here"
									className="w-full"
									type="text"
									// defaultValue={messageContent || undefined}
									value={messageContent || ''}
									onChange={e => {
										setMessageContent(() => e.target.value)
									}}
								/>

								{isFetchingSuggestions || isRefetchingSuggestions ? (
									<div className="rotate size-6 animate-spin rounded-full border-4 border-solid  border-l-primary" />
								) : (
									<div
										className="flex size-6 shrink-0 cursor-pointer items-center justify-center rounded-full bg-background ring-1 ring-border"
										onClick={() => {
											refetchSuggestions().catch(error =>
												console.error(error)
											)
										}}
									>
										<div className="translate-y-px">
											<SparklesIcon size={12} />
										</div>
									</div>
								)}

								<Button type="submit" className="rounded-full" disabled={isBusy}>
									<SendIcon className="size-4" />
								</Button>
							</form>
						</CardFooter>
					</div>
				</>
			) : (
				<div className="flex h-full flex-col items-center justify-center bg-[#ebe5de] dark:bg-[#111b21]">
					<Icons.pointer className="size-4" />
					<p className="text-lg font-semibold">No Conversation Selected</p>
					<p className="text-sm">Select a conversation from the side list</p>
				</div>
			)}
		</div>
	)
}

export default ChatCanvas
