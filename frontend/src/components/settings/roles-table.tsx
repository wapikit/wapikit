'use client'

import { useSearchParams } from 'next/navigation'
import { TableComponent } from '../tables/table'
import { type OrganizationMemberSchema, useGetOrganizationMembers } from 'root/.generated'
import { RolesTableColumns } from '../tables/columns'

const TeamTable = () => {
	const searchParams = useSearchParams()

	const page = Number(searchParams.get('page') || 1)
	const pageLimit = Number(searchParams.get('limit') || 0) || 10
	const teamMemberResponse = useGetOrganizationMembers({})
	const totalUsers = teamMemberResponse.data?.paginationMeta?.total || 0
	const pageCount = Math.ceil(totalUsers / pageLimit)
	const teamMembers: OrganizationMemberSchema[] = teamMemberResponse.data?.members || []

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
