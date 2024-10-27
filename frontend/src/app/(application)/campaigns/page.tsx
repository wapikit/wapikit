'use client'

import BreadCrumb from '~/components/breadcrumb'
import { CampaignTableColumns } from '~/components/tables/columns'
import { TableComponent } from '~/components/tables/table'
import { buttonVariants } from '~/components/ui/button'
import { Heading } from '~/components/ui/heading'
import { Separator } from '~/components/ui/separator'
import { useDeleteCampaignById, useGetCampaigns, type CampaignSchema } from 'root/.generated'
import { Plus } from 'lucide-react'
import Link from 'next/link'
import { clsx } from 'clsx'
import { useRouter, useSearchParams } from 'next/navigation'
import { useState } from 'react'
import { errorNotification, materialConfirm, successNotification } from '~/reusable-functions'

const breadcrumbItems = [{ title: 'campaigns', link: '/campaigns' }]

const CampaignsPage = () => {
	// ! TODO: Implement CampaignsPage
	// * 1. Create a table of campaigns with pagination enabled
	// * 2. Handle query params for pagination
	// * 3. List actions: Delete, Export, Create a new campaign

	const searchParams = useSearchParams()
	const router = useRouter()
	const deleteCampaignMutation = useDeleteCampaignById()

	const page = Number(searchParams.get('page') || 1)
	const pageLimit = Number(searchParams.get('limit') || 0) || 10
	const [order] = useState()
	const [status] = useState()

	const contactResponse = useGetCampaigns({
		per_page: pageLimit || 10,
		page: page || 1,
		...(order ? { order: order } : {}),
		...(status ? { status: status } : {})
	})

	const totalCampaigns = contactResponse.data?.paginationMeta?.total || 0
	const pageCount = Math.ceil(totalCampaigns / pageLimit)
	const campaigns: CampaignSchema[] = contactResponse.data?.campaigns || []

	async function deleteCampaign(campaignId: string) {
		try {
			if (!campaignId) return

			const confirmation = await materialConfirm({
				title: 'Delete Campaign',
				description: 'Are you sure you want to delete this campaign?'
			})

			if (!confirmation) return

			const { data } = await deleteCampaignMutation.mutateAsync({
				id: campaignId
			})

			if (data) {
				// show success notification
				successNotification({
					message: 'Campaign deleted successfully'
				})
			} else {
				// show error notification
				errorNotification({
					message: 'Failed to delete campaign'
				})
			}
		} catch (error) {
			console.error('Error deleting campaign', error)
			errorNotification({
				message: 'Error deleting campaign'
			})
		}
	}

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
						href={'/campaigns/new-or-edit'}
						className={clsx(buttonVariants({ variant: 'default' }))}
					>
						<Plus className="mr-2 h-4 w-4" /> Add New
					</Link>
				</div>
				<Separator />

				<TableComponent
					searchKey="name"
					pageNo={page}
					columns={CampaignTableColumns}
					totalUsers={totalCampaigns}
					data={campaigns}
					pageCount={pageCount}
					actions={[
						{
							icon: 'edit',
							label: 'Edit',
							onClick: (campaignId: string) => {
								// redirect to the edit page with id in search param
								router.push(`/campaigns/new-or-edit?id=${campaignId}`)
							}
						},
						{
							icon: 'trash',
							label: 'Delete',
							onClick: (campaignId: string) => {
								deleteCampaign(campaignId).catch(console.error)
							}
						}
					]}
				/>
			</div>
		</>
	)
}

export default CampaignsPage
