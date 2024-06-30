'use client'

import { useSearchParams } from 'next/navigation'
import { TableComponent } from '../tables/table'
import { useGetOrganizationRoles, type OrganizationRoleSchema } from 'root/.generated'
import { RolesTableColumns } from '../tables/columns'

const TeamTable = () => {
	const searchParams = useSearchParams()
	const page = Number(searchParams.get('page') || 1)
	const pageLimit = Number(searchParams.get('limit') || 0) || 10
	const rolesResponse = useGetOrganizationRoles({
		page: page || 1,
		per_page: pageLimit || 10
	})
	const totalUsers = rolesResponse.data?.paginationMeta?.total || 0
	const pageCount = Math.ceil(totalUsers / pageLimit)
	const teamMembers: OrganizationRoleSchema[] = rolesResponse.data?.roles || []

	return (
		<TableComponent
			searchKey="name"
			pageNo={page}
			columns={RolesTableColumns}
			totalUsers={totalUsers}
			data={teamMembers}
			pageCount={pageCount}
		/>
	)
}

export default TeamTable
