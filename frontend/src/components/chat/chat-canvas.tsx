'use client'

import { ScrollArea } from '~/components/ui/scroll-area'
import { CardHeader, CardContent, CardFooter } from '../ui/card'
import { Separator } from '../ui/separator'
import { Input } from '../ui/input'
import { Button } from '../ui/button'
import { SendIcon, Image as ImageIcon, MoreVerticalIcon } from 'lucide-react'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger
} from '../ui/dropdown-menu'
import { Icons } from '../icons'
import { useState } from 'react'
import {
	useAssignConversation,
	useGetOrganizationMembers,
	useUnassignConversation
} from 'root/.generated'
import MessageRenderer from './message-renderer'
import { useRouter } from 'next/navigation'
import { errorNotification, successNotification } from '~/reusable-functions'
import { Modal } from '../ui/modal'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '../ui/form'
import { ReloadIcon } from '@radix-ui/react-icons'
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

const ChatCanvas = () => {
	const [isBusy, setIsBusy] = useState(false)
	const [isConversationAssignModalOpen, setIsConversationAssignModalOpen] = useState(false)
	const { currentConversation } = useConversationInboxStore()

	const router = useRouter()
	const { writeProperty } = useLayoutStore()

	const assignConversationMutation = useAssignConversation()
	const unassignConversationMutation = useUnassignConversation()

	const assignConversationForm = useForm<z.infer<typeof AssignConversationForm>>({
		resolver: zodResolver(AssignConversationForm)
	})

	const { data: organizationMembersResponse, refetch: refetchMembers } =
		useGetOrganizationMembers({
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
			label: 'Assign to',
			icon: 'user',
			onClick() {
				setIsConversationAssignModalOpen(true)
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
					contactSheetData: {
						name: '',
						attributes: {},
						createdAt: '',
						lists: [],
						phone: '',
						uniqueId: ''
					}
				})
			}
		}
	]

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
												Message Template
												<Button
													disabled={isBusy}
													size={'sm'}
													variant={'secondary'}
													type="button"
													onClick={e => {
														e.preventDefault()
														refetchMembers().catch(error =>
															console.error(error)
														)
													}}
												>
													<ReloadIcon className="size-3" />
												</Button>
											</FormLabel>
											<FormControl>
												<Select
													disabled={isBusy}
													onValueChange={e => {
														field.onChange(e)
													}}
													name="templateId"
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
																'Select message template'}
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
								Invite Now
							</Button>
						</form>
					</Form>
				</div>
			</Modal>

			{currentConversation ? (
				<>
					<CardHeader className="item-center flex !flex-row justify-between rounded-t-md  bg-primary p-3 py-2">
						<div className="flex flex-row gap-3 ">
							<Image
								src={'/assets/empty-pfp.png'}
								height={50}
								width={50}
								className="object-fit aspect-square h-10 w-10 rounded-full"
								alt={`${currentConversation.uniqueId} avatar`}
							/>
							<div className="flex flex-col">
								<p className="text-base text-primary-foreground">
									{currentConversation.contact.name}
								</p>
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
						<CardContent className="relative h-full w-full  bg-[#ebe5de] !py-4">
							<div className='absolute inset-0 z-20 h-full w-full  bg-[url("/assets/conversations-canvas-bg.png")] bg-repeat opacity-20 ' />
							{currentConversation.messages.map((message, index) => {
								return (
									<div className="relative z-30 w-full " key={index}>
										<MessageRenderer
											message={message}
											isActionsEnabled={true}
										/>
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
				</>
			) : (
				<div className="flex h-full flex-col items-center justify-center bg-[#ebe5de]">
					<Icons.pointer className="size-4" />
					<p className="text-lg font-semibold">No Conversation Selected</p>
					<p className="text-sm">Select a conversation from the side list</p>
				</div>
			)}
		</div>
	)
}

export default ChatCanvas
