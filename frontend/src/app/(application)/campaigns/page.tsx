'use client'

import BreadCrumb from '~/components/breadcrumb'
import { CampaignTableColumns } from '~/components/tables/columns'
import { TableComponent } from '~/components/tables/table'
import { buttonVariants } from '~/components/ui/button'
import { Heading } from '~/components/ui/heading'
import { Separator } from '~/components/ui/separator'
import { useGetCampaigns, type CampaignSchema } from 'root/.generated'
import { Plus } from 'lucide-react'
import Link from 'next/link'
import { clsx } from 'clsx'
import { useSearchParams } from 'next/navigation'

const breadcrumbItems = [{ title: 'campaigns', link: '/campaigns' }]

const CampaignsPage = () => {
	// ! TODO: Implement CampaignsPage
	// * 1. Create a table of campaigns with pagination enabled
	// * 2. Handle query params for pagination
	// * 3. List actions: Delete, Export, Create a new campaign

	const searchParams = useSearchParams()

	const page = Number(searchParams.get('page') || 1)
	const pageLimit = Number(searchParams.get('limit') || 0) || 10
	// const offset = (page - 1) * pageLimit

	const contactResponse = useGetCampaigns({})

	const totalCampaigns = contactResponse.data?.paginationMeta?.total || 0
	const pageCount = Math.ceil(totalCampaigns / pageLimit)
	const campaigns: CampaignSchema[] = contactResponse.data?.campaigns || []

	return (
		<>
			<div className="flex-1 space-y-4  p-4 pt-6 md:p-8">
				<BreadCrumb items={breadcrumbItems} />
				<div className="flex items-start justify-between">
					<Heading
						title={`Campaigns (${totalCampaigns})`}
						description="Manage campaigns"
					/>

					<Link
						href={'/campaigns/new'}
						className={clsx(buttonVariants({ variant: 'default' }))}
					>
						<Plus className="mr-2 h-4 w-4" /> Add New
					</Link>
				</div>
				<Separator />

				<TableComponent
					searchKey="country"
					pageNo={page}
					columns={CampaignTableColumns}
					totalUsers={totalCampaigns}
					data={campaigns}
					pageCount={pageCount}
				/>
			</div>
		</>
	)
}

export default CampaignsPage
