'use client'

import { useSearchParams } from 'next/navigation'
import { TableComponent } from '../tables/table'
import {
	type GetOrganizationRolesResponseSchema,
	useDeleteOrganizationRoleById,
	type OrganizationRoleSchema
} from 'root/.generated'
import { RolesTableColumns } from '../tables/columns'
import { errorNotification, materialConfirm, successNotification } from '~/reusable-functions'
import { type Dispatch, type SetStateAction } from 'react'

const RolesTable: React.FC<{
	setRoleToEditId: Dispatch<SetStateAction<string | null>>
	rolesResponse?: GetOrganizationRolesResponseSchema
}> = ({ setRoleToEditId, rolesResponse }) => {
	const searchParams = useSearchParams()
	const deleteRoleMutation = useDeleteOrganizationRoleById()
	const page = Number(searchParams.get('page') || 1)
	const pageLimit = Number(searchParams.get('limit') || 0) || 10

	const totalUsers = rolesResponse?.paginationMeta?.total || 0
	const pageCount = Math.ceil(totalUsers / pageLimit)
	const roles: OrganizationRoleSchema[] = rolesResponse?.roles || []

	async function handleDeleteRole(roleId: string) {
		try {
			if (!roleId) return

			const confirmation = await materialConfirm({
				title: 'Delete Role',
				description: 'Are you sure you want to delete this role?'
			})

			if (!confirmation) return

			const { data } = await deleteRoleMutation.mutateAsync({
				id: roleId
			})

			if (data) {
				successNotification({
					message: 'Role deleted successfully'
				})
			} else {
				errorNotification({
					message: 'Error deleting role'
				})
			}
		} catch (error) {
			console.error('Error deleting role', error)
			errorNotification({
				message: 'Error deleting role'
			})
		}
	}

	return (
		<TableComponent
			searchKey="name"
			pageNo={page}
			columns={RolesTableColumns}
			totalUsers={totalUsers}
			data={roles}
			pageCount={pageCount}
			actions={[
				{
					label: 'Edit',
					onClick: (roleId: string) => {
						setRoleToEditId(() => roleId)
					},
					icon: 'edit'
				},
				{
					label: 'Delete',
					onClick: (roleId: string) => {
						handleDeleteRole(roleId).catch(console.error)
					},
					icon: 'trash'
				}
			]}
		/>
	)
}

export default RolesTable
