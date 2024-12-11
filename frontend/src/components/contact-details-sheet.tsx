'use client'

import { useLayoutStore } from '~/store/layout.store'
import { Label } from './ui/label'
import { Sheet, SheetContent, SheetHeader, SheetTitle } from './ui/sheet'
import { Input } from './ui/input'
import { usePathname, useRouter } from 'next/navigation'
import { Icons } from './icons'
import { Badge } from './ui/badge'
import dayjs from 'dayjs'
import { Button } from './ui/button'
import { errorNotification, materialConfirm, successNotification } from '~/reusable-functions'
import { useDeleteContactById } from 'root/.generated'

const ContactDetailsSheet = () => {
	const { writeProperty, contactSheetData } = useLayoutStore()
	const pathname = usePathname()
	const router = useRouter()

	const deleteContactByIdMutation = useDeleteContactById()

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
			open={!!contactSheetData}
			onOpenChange={isOpen => {
				if (!isOpen) {
					writeProperty({ contactSheetData: null })
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
				<div className="relative flex h-screen flex-col justify-between gap-6">
					{contactSheetData ? (
						<div className="flex flex-col gap-3 text-base">
							{Object.keys(contactSheetData || {})
								.filter(key => key !== 'uniqueId')
								.sort((a, b) => sortedKeys.indexOf(a) - sortedKeys.indexOf(b))
								.map(key => {
									let IconToRender = Icons.user

									if (key === 'createdAt') {
										IconToRender = Icons.calendar
									} else if (key === 'phone') {
										IconToRender = Icons.phone
									} else if (key === 'attributes') {
										IconToRender = Icons.jsonBrackets
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
														{contactSheetData.lists.length === 0 && (
															<Badge variant={'outline'}>None</Badge>
														)}
														{contactSheetData.lists.map(
															(list, index) => {
																if (index > 2) {
																	return null
																}
																return (
																	<Badge key={list.uniqueId}>
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
										<div className="flex flex-row items-center gap-2" key={key}>
											<Label
												htmlFor={key}
												className="flex items-center gap-2 text-left"
											>
												<IconToRender className="size-4" />
												{key}:
											</Label>
											<span className="text-sm font-semibold">
												{key === 'attributes'
													? JSON.stringify(contactSheetData['attributes'])
													: key === 'createdAt'
														? dayjs(contactSheetData.createdAt).format(
																'DD MMM, YYYY'
															)
														: contactSheetData[key] || 'N/A'}
											</span>
										</div>
									)
								})}
						</div>
					) : null}

					{contactSheetData ? (
						<div className="sticky bottom-0 flex w-full flex-row items-center justify-between gap-3">
							<Button
								onClick={() => {
									router.push(
										`/contacts/new-or-edit?id=${contactSheetData.uniqueId}`
									)
								}}
								className="flex w-full flex-row gap-2"
								variant={'secondary'}
							>
								<Icons.edit className="size-4" />
								Edit
							</Button>
							<Button variant={'secondary'} className="flex w-full flex-row gap-2">
								<Icons.xCircle className="size-4" />
								Block
							</Button>
							<Button
								variant={'destructive'}
								onClick={() => {
									deleteContact(contactSheetData.uniqueId).catch(error =>
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
				</div>
			</SheetContent>
		</Sheet>
	)
}

export default ContactDetailsSheet
