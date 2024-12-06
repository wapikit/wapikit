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
	UserPermissionLevel,
	useUpdateOrganizationMemberRoleById,
	useGetOrganizationRoles
} from 'root/.generated'
import { Plus } from 'lucide-react'
import { clsx } from 'clsx'
import { useSearchParams } from 'next/navigation'
import { Modal } from '~/components/ui/modal'
import { useEffect, useMemo, useState } from 'react'
import { errorNotification, materialConfirm, successNotification } from '~/reusable-functions'
import { Input } from '~/components/ui/input'
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
import { NewTeamMemberInviteFormSchema, UpdateOrganizationMemberRolesFormSchema } from '~/schema'
import type { z } from 'zod'
import { MultiSelect } from '~/components/multi-select'

const breadcrumbItems = [{ title: 'Members', link: '/members' }]

const MembersPage = () => {
	const searchParams = useSearchParams()

	const [isInvitationModalOpen, setIsInvitationModalOpen] = useState(false)
	const [memberToEditId, setMemberToEditId] = useState<string | null>(null)
	const [isBusy, setIsBusy] = useState(false)

	const page = Number(searchParams.get('page') || 1)
	const pageLimit = Number(searchParams.get('limit') || 0) || 10
	const sortBy = searchParams.get('sortOrder')

	const { data: membersResponse, refetch: refetchMembers } = useGetOrganizationMembers({
		page: page || 1,
		per_page: pageLimit || 10,
		sortBy: sortBy ? (sortBy as OrderEnum) : undefined
	})

	const { data: allRoles } = useGetOrganizationRoles({
		per_page: 50,
		page: 1,
		sortBy: 'asc'
	})

	const newMemberInviteForm = useForm<z.infer<typeof NewTeamMemberInviteFormSchema>>({
		resolver: zodResolver(NewTeamMemberInviteFormSchema),
		defaultValues: {
			email: ''
		}
	})

	const memberUpdateForm = useForm<z.infer<typeof UpdateOrganizationMemberRolesFormSchema>>({
		resolver: zodResolver(UpdateOrganizationMemberRolesFormSchema),
		defaultValues: {
			roles: memberToEditId
				? membersResponse?.members
						.find(member => member.uniqueId === memberToEditId)
						?.roles?.map(role => role.uniqueId)
				: []
		}
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
	const updateMemberRolesMutation = useUpdateOrganizationMemberRoleById()

	async function inviteUser() {
		try {
			console.log(newMemberInviteForm.getValues())
			setIsBusy(true)
			const confirmation = await materialConfirm({
				description: 'Are you sure you want to invite this user?',
				title: 'Invite User'
			})

			if (!confirmation) return

			const response = await inviteUserMutation.mutateAsync({
				data: {
					email: newMemberInviteForm.getValues('email'),
					accessLevel: UserPermissionLevel.Member
				}
			})

			if (response.invite) {
				successNotification({
					message: 'User invited successfully.'
				})
				newMemberInviteForm.reset()
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

	async function updateMemberRoles(
		data: z.infer<typeof UpdateOrganizationMemberRolesFormSchema>
	) {
		try {
			if (!memberToEditId) return

			const response = await updateMemberRolesMutation.mutateAsync({
				id: memberToEditId,
				data: {
					updatedRoleIds: data.roles
				}
			})

			if (response.isRoleUpdated) {
				successNotification({
					message: 'Member roles updated successfully.'
				})
				memberUpdateForm.reset()
				setMemberToEditId(null)
				await refetchMembers()
			} else {
				errorNotification({
					message: 'Something went wrong, While updating member roles. Please try again.'
				})
			}
		} catch (error) {
			console.error(error)
			errorNotification({
				message: 'Something went wrong, While updating member roles. Please try again.'
			})
		}
	}

	useEffect(() => {
		if (memberUpdateForm.formState.isDirty) return

		if (memberToEditId) {
			memberUpdateForm.setValue(
				'roles',
				membersResponse?.members
					.find(member => member.uniqueId === memberToEditId)
					?.roles?.map(role => role.uniqueId) || [],
				{
					shouldValidate: true,
					shouldDirty: true
				}
			)
		}
	}, [memberToEditId, memberUpdateForm, membersResponse?.members])

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
					<Form {...newMemberInviteForm}>
						<form
							onSubmit={newMemberInviteForm.handleSubmit(inviteUser)}
							className="w-full space-y-8"
						>
							<div className="flex flex-col gap-8">
								<FormField
									control={newMemberInviteForm.control}
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
							</div>
							<Button disabled={isBusy} className="ml-auto mr-0 w-full" type="submit">
								Invite Now
							</Button>
						</form>
					</Form>
				</div>
			</Modal>

			{/* edit member details modal */}
			<Modal
				title="Edit Organization Member"
				description="update the member roles."
				isOpen={!!memberToEditId}
				onClose={() => {
					setMemberToEditId(null)
				}}
			>
				<div className="flex w-full items-center justify-end space-x-2 pt-6">
					<Form {...memberUpdateForm}>
						<form
							onSubmit={memberUpdateForm.handleSubmit(updateMemberRoles)}
							className="w-full space-y-8"
						>
							<div className="flex flex-col gap-8">
								<FormField
									control={memberUpdateForm.control}
									name="roles"
									render={({}) => (
										<FormItem className="tablet:w-3/4 tablet:gap-2 desktop:w-1/2 flex flex-col gap-1 ">
											<FormLabel>Select the permissions</FormLabel>
											<MultiSelect
												options={(allRoles?.roles || []).map(role => {
													return {
														value: role.uniqueId,
														label: role.name
													}
												})}
												onValueChange={e => {
													console.log({ e })
													memberUpdateForm.setValue(
														'roles',
														e as string[],
														{
															shouldValidate: true
														}
													)
												}}
												defaultValue={memberUpdateForm.watch('roles')}
												placeholder="Select permissions"
												variant="default"
											/>
											<FormMessage />
										</FormItem>
									)}
								/>
							</div>
							<Button disabled={isBusy} className="ml-auto mr-0 w-full" type="submit">
								Update Roles
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
							label: 'Update Roles',
							onClick: (memberId: string) => {
								setMemberToEditId(() => memberId)
							}
						}
					]}
				/>
			</div>
		</>
	)
}

export default MembersPage
