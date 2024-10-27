'use client'

import BreadCrumb from '~/components/breadcrumb'
import { OrganizationMembersTableColumns } from '~/components/tables/columns'
import { TableComponent } from '~/components/tables/table'
import { Button, buttonVariants } from '~/components/ui/button'
import { Heading } from '~/components/ui/heading'
import { Separator } from '~/components/ui/separator'
import {
	type OrderEnum,
	useGetOrganizationMembers,
	useCreateOrganizationInvite,
	UserPermissionLevel
} from 'root/.generated'
import { Plus } from 'lucide-react'
import { clsx } from 'clsx'
import { useRouter, useSearchParams } from 'next/navigation'
import { Modal } from '~/components/ui/modal'
import { useMemo, useState } from 'react'
import { errorNotification, materialConfirm, successNotification } from '~/reusable-functions'
import { Input } from '~/components/ui/input'
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue
} from '~/components/ui/select'
import { listStringEnumMembers } from 'ts-enum-utils'
import {
	Form,
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage
} from '~/components/ui/form'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { NewTeamMemberInviteFormSchema } from '~/schema'
import type { z } from 'zod'

const breadcrumbItems = [{ title: 'Members', link: '/members' }]

const MembersPage = () => {
	const searchParams = useSearchParams()
	const router = useRouter()

	const [isInvitationModalOpen, setIsInvitationModalOpen] = useState(false)
	const [isBusy, setIsBusy] = useState(false)

	const form = useForm<z.infer<typeof NewTeamMemberInviteFormSchema>>({
		resolver: zodResolver(NewTeamMemberInviteFormSchema),
		defaultValues: {
			email: '',
			accessLevel: UserPermissionLevel.Member
		}
	})

	const page = Number(searchParams.get('page') || 1)
	const pageLimit = Number(searchParams.get('limit') || 0) || 10
	const sortBy = searchParams.get('sortOrder')

	const { data: membersResponse, refetch: refetchMembers } = useGetOrganizationMembers({
		page: page || 1,
		per_page: pageLimit || 10,
		sortBy: sortBy ? (sortBy as OrderEnum) : undefined
	})

	const organizationMembersList = useMemo(() => {
		return membersResponse?.members || []
	}, [membersResponse])

	const paginationMeta = useMemo(() => {
		return membersResponse?.paginationMeta
	}, [membersResponse])

	const totalUsers = paginationMeta?.total || 0
	const pageCount = Math.ceil(totalUsers / pageLimit)

	const inviteUserMutation = useCreateOrganizationInvite()

	async function inviteUser() {
		try {
			console.log(form.getValues())
			setIsBusy(true)
			const confirmation = await materialConfirm({
				description: 'Are you sure you want to invite this user?',
				title: 'Invite User'
			})

			if (!confirmation) return

			const response = await inviteUserMutation.mutateAsync({
				data: {
					accessLevel: form.getValues('accessLevel'),
					email: form.getValues('email')
				}
			})

			console.log(response)

			if (response.invite) {
				successNotification({
					message: 'User invited successfully.'
				})
				form.reset()
				setIsInvitationModalOpen(false)
				await refetchMembers()
			} else {
				errorNotification({
					message: 'Something went wrong, While inviting a user. Please try again.'
				})
			}
		} catch (error) {
			console.error(error)
			errorNotification({
				message: 'Something went wrong, While inviting a user. Please try again.'
			})
		} finally {
			setIsBusy(false)
		}
	}

	return (
		<>
			{/* invitation form modal */}
			<Modal
				title="Invite Team Member"
				description="an email would be sent to them."
				isOpen={isInvitationModalOpen}
				onClose={() => {
					setIsInvitationModalOpen(false)
				}}
			>
				<div className="flex w-full items-center justify-end space-x-2 pt-6">
					<Form {...form}>
						<form onSubmit={form.handleSubmit(inviteUser)} className="w-full space-y-8">
							<div className="flex flex-col gap-8">
								<FormField
									control={form.control}
									name="email"
									render={({ field }) => (
										<FormItem>
											<FormLabel>Email</FormLabel>
											<FormControl>
												<Input
													disabled={isBusy}
													placeholder="Email"
													{...field}
													autoComplete="off"
												/>
											</FormControl>
											<FormMessage />
										</FormItem>
									)}
								/>

								<FormField
									control={form.control}
									name="accessLevel"
									render={({ field }) => (
										<FormItem>
											<FormLabel>Access Level</FormLabel>
											<Select
												disabled={isBusy}
												onValueChange={field.onChange}
												value={field.value}
												// defaultValue={field.value}
											>
												<FormControl>
													<SelectTrigger>
														<SelectValue
															defaultValue={field.value}
															placeholder="Select Access Level"
														/>
													</SelectTrigger>
												</FormControl>
												<SelectContent>
													{listStringEnumMembers(UserPermissionLevel).map(
														status => {
															return (
																<SelectItem
																	key={status.name}
																	value={status.value}
																>
																	{status.name}
																</SelectItem>
															)
														}
													)}
												</SelectContent>
											</Select>
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

			<div className="flex-1 space-y-4  p-4 pt-6 md:p-8">
				<BreadCrumb items={breadcrumbItems} />

				<div className="flex items-start justify-between">
					<Heading title={`Team Members (${totalUsers})`} description="Manage members" />
					<Button
						onClick={() => {
							setIsInvitationModalOpen(true)
						}}
						className={clsx(buttonVariants({ variant: 'default' }))}
					>
						<Plus className="mr-2 h-4 w-4" /> Add New
					</Button>
				</div>
				<Separator />

				<TableComponent
					searchKey="name"
					pageNo={page}
					columns={OrganizationMembersTableColumns}
					totalUsers={totalUsers}
					data={organizationMembersList}
					pageCount={pageCount}
					actions={[
						{
							icon: 'edit',
							label: 'Edit',
							onClick: (contactId: string) => {
								// redirect to the edit page with id in search param
								router.push(`/contacts/new-or-edit?id=${contactId}`)
							}
						}
					]}
				/>
			</div>
		</>
	)
}

export default MembersPage
