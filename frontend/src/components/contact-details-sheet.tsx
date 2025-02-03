'use client'

import { useLayoutStore } from '~/store/layout.store'
import { Label } from './ui/label'
import { Sheet, SheetContent, SheetHeader, SheetTitle } from './ui/sheet'
import { usePathname, useRouter } from 'next/navigation'
import { Icons } from './icons'
import { Badge } from './ui/badge'
import dayjs from 'dayjs'
import { Button } from '~/components/ui/button'
import { errorNotification, materialConfirm, successNotification } from '~/reusable-functions'
import { useDeleteContactById, useGetContactById } from 'root/.generated'
import { Card, CardContent, CardTitle } from './ui/card'
import Link from 'next/link'
import LoadingSpinner from './loader'

const ContactDetailsSheet = () => {
	const { writeProperty, contactSheetContactId } = useLayoutStore()
	const pathname = usePathname()
	const router = useRouter()

	const deleteContactByIdMutation = useDeleteContactById()

	const { data, isFetching } = useGetContactById(contactSheetContactId || '')

	const contactData = data?.contact

	async function deleteContact(contactId: string) {
		try {
			if (!contactId) return

			const confirmation = await materialConfirm({
				title: 'Delete Contact',
				description: 'Are you sure you want to delete this contact?'
			})

			if (!confirmation) return

			const { data } = await deleteContactByIdMutation.mutateAsync({
				id: contactId
			})

			if (data) {
				successNotification({
					message: 'Contact deleted successfully'
				})

				router.push('/contacts')
			} else {
				errorNotification({
					message: 'Failed to delete contact'
				})
			}
		} catch (error) {
			console.error('Error deleting contact', error)
			errorNotification({
				message: 'Error deleting contact'
			})
		}
	}

	const sortedKeys = ['name', 'phone', 'createdAt', 'lists', 'attributes']

	return (
		<Sheet
			open={!!contactData}
			onOpenChange={isOpen => {
				if (!isOpen) {
					writeProperty({ contactSheetContactId: null })
					if (pathname === '/contacts') {
						router.push('/contacts')
					}
				}
			}}
		>
			<SheetContent className="flex h-screen flex-col items-start justify-start gap-6">
				<SheetHeader>
					<SheetTitle>Contact Info</SheetTitle>
				</SheetHeader>

				{isFetching ? (
					<>
						<LoadingSpinner />
					</>
				) : (
					<>
						<Card className="h-fit w-full py-4">
							<CardContent className="flex h-fit flex-col justify-start gap-6">
								{contactData ? (
									<div className="flex flex-col gap-3 text-base">
										{Object.keys(contactData || {})
											.filter(key => key !== 'uniqueId')
											.sort(
												(a, b) =>
													sortedKeys.indexOf(a) - sortedKeys.indexOf(b)
											)
											.map(k => {
												const key: keyof typeof contactData =
													k as keyof typeof contactData

												if (key === 'conversations') return null
												let IconToRender: any = Icons.user

												if (key === 'createdAt') {
													IconToRender = Icons.calendar
												} else if (key === 'phone') {
													IconToRender = Icons.phone
												} else if (key === 'attributes') {
													IconToRender = Icons.code
												}

												if (key === 'lists') {
													IconToRender = Icons.dashboard

													return (
														<div
															className="flex flex-row items-center gap-2"
															key={key}
														>
															<Label
																htmlFor={key}
																className="flex items-center gap-2 text-left"
															>
																<IconToRender className="size-4" />
																{key}:
															</Label>
															<span className="text-base font-semibold">
																<div className="flex flex-wrap items-center justify-center gap-0.5 truncate">
																	{contactData.lists.length ===
																		0 && (
																		<Badge variant={'outline'}>
																			None
																		</Badge>
																	)}
																	{contactData.lists.map(
																		(list, index) => {
																			if (index > 2) {
																				return null
																			}
																			return (
																				<Badge
																					key={
																						list.uniqueId
																					}
																				>
																					{list.name}
																				</Badge>
																			)
																		}
																	)}
																</div>
															</span>
														</div>
													)
												}

												// ! TODO: handle conversation section here
												// ! TODO: add tags rendering here

												return (
													<div
														className="flex flex-row items-center gap-2"
														key={key}
													>
														<Label
															htmlFor={key}
															className="flex items-center gap-2 text-left"
														>
															<IconToRender className="size-4" />
															{key}:
														</Label>
														<span className="text-sm font-semibold">
															{key === 'attributes'
																? JSON.stringify(
																		contactData['attributes']
																	)
																: key === 'createdAt'
																	? dayjs(
																			contactData.createdAt
																		).format('DD MMM, YYYY')
																	: contactData[key] || 'N/A'}
														</span>
													</div>
												)
											})}
									</div>
								) : null}
							</CardContent>
						</Card>

						{contactData?.conversations?.length ? (
							<Card className="flex w-full flex-col gap-4 py-4">
								<CardTitle className="flex flex-row items-center justify-start gap-2 px-4">
									<Icons.message className="size-4" />
									Conversations
								</CardTitle>
								<CardContent className="flex w-full flex-col gap-2 divide-y">
									{contactData.conversations.map((conversation, index) => (
										<Link
											key={index}
											href={`/conversations?id=${conversation.uniqueId}`} // Use a unique ID if available
											className="block w-full rounded-md bg-gray-50 px-2 py-3 transition-all hover:bg-gray-100"
										>
											<div className="flex flex-col gap-2">
												{/* Conversation Date */}
												<div className="flex flex-row items-center gap-2">
													<Icons.calendar className="size-4 text-gray-500" />
													<span className="text-sm font-semibold">
														{dayjs(conversation.createdAt).format(
															'DD MMM, YYYY'
														)}
													</span>
												</div>

												{/* First Message Snippet */}
												<div className="flex flex-row items-center gap-2">
													<Icons.message className="size-4 text-gray-500" />
													<span className="max-w-[300px] truncate text-sm text-gray-700">
														{conversation.messages?.length
															? conversation.messages[0]?.messageData
																	?.text + '...'
															: 'No messages yet'}
													</span>
												</div>

												{/* ! TODO: fetch the campaign name and assigned to member here */}
												{/* Show Campaign Name if Exists */}
												{/* {conversation.campaignId && (
												<div className="flex flex-row items-center gap-2 text-xs text-gray-600">
													<Icons.rocket className="size-4 text-gray-500" />
													Campaign:{' '}
													<span className="font-medium">
														{conversation.campaignId}
													</span>
												</div>
											)} */}
											</div>
										</Link>
									))}
								</CardContent>
							</Card>
						) : null}

						{contactData ? (
							<div className="sticky bottom-0 flex w-full flex-row items-center justify-between gap-3">
								<Button
									onClick={() => {
										router.push(
											`/contacts/new-or-edit?id=${contactData.uniqueId}`
										)
									}}
									className="flex w-full flex-row gap-2"
									variant={'secondary'}
								>
									<Icons.edit className="size-4" />
									Edit
								</Button>
								<Button
									variant={'secondary'}
									className="flex w-full flex-row gap-2"
								>
									<Icons.xCircle className="size-4" />
									Block
								</Button>
								<Button
									variant={'destructive'}
									onClick={() => {
										deleteContact(contactData.uniqueId).catch(error =>
											console.error(error)
										)
									}}
									className="flex w-full flex-row gap-2"
								>
									<Icons.trash className="size-4" />
									Delete
								</Button>
							</div>
						) : null}
					</>
				)}
			</SheetContent>
		</Sheet>
	)
}

export default ContactDetailsSheet
