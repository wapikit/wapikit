'use client'

import BreadCrumb from '~/components/breadcrumb'
import { OrgnizationMembersTableColumns } from '~/components/tables/columns'
import { TableComponent } from '~/components/tables/table'
import { buttonVariants } from '~/components/ui/button'
import { Heading } from '~/components/ui/heading'
import { Separator } from '~/components/ui/separator'
import { type OrganizationMemberSchema, useGetOrganizationMembers } from 'root/.generated'
import { Plus } from 'lucide-react'
import Link from 'next/link'
import { clsx } from 'clsx'
import { useSearchParams } from 'next/navigation'

const breadcrumbItems = [{ title: 'Members', link: '/members' }]

const MembersPage = () => {
	// ! TODO:

	const searchParams = useSearchParams()

	const page = Number(searchParams.get('page') || 1)
	const pageLimit = Number(searchParams.get('limit') || 0) || 10
	// const offset = (page - 1) * pageLimit

	const membersResponse = useGetOrganizationMembers({})

	const totalUsers = membersResponse.data?.paginationMeta?.total || 0
	const pageCount = Math.ceil(totalUsers / pageLimit)
	const contacts: OrganizationMemberSchema[] = membersResponse.data?.members || []

	return (
		<>
			<div className="flex-1 space-y-4  p-4 pt-6 md:p-8">
				<BreadCrumb items={breadcrumbItems} />

				<div className="flex items-start justify-between">
					<Heading title={`Team Members (${totalUsers})`} description="Manage members" />
					<Link
						href={'/members/new'}
						className={clsx(buttonVariants({ variant: 'default' }))}
					>
						<Plus className="mr-2 h-4 w-4" /> Add New
					</Link>
				</div>
				<Separator />

				<TableComponent
					searchKey="country"
					pageNo={page}
					columns={OrgnizationMembersTableColumns}
					totalUsers={totalUsers}
					data={contacts}
					pageCount={pageCount}
				/>
			</div>
		</>
	)
}

export default MembersPage
