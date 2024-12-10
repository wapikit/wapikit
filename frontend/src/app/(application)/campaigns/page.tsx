'use client'

import BreadCrumb from '~/components/breadcrumb'
import { CampaignTableColumns } from '~/components/tables/columns'
import { TableComponent } from '~/components/tables/table'
import { buttonVariants } from '~/components/ui/button'
import { Heading } from '~/components/ui/heading'
import { Separator } from '~/components/ui/separator'
import {
	useDeleteCampaignById,
	useGetCampaignAnalyticsById,
	useGetCampaignById,
	useGetCampaigns,
	useUpdateCampaignById,
	type CampaignSchema
} from 'root/.generated'
import { Plus } from 'lucide-react'
import Link from 'next/link'
import { clsx } from 'clsx'
import { notFound, useRouter, useSearchParams } from 'next/navigation'
import { useState } from 'react'
import { errorNotification, materialConfirm, successNotification } from '~/reusable-functions'
import type { TableCellActionProps } from '~/types'

const breadcrumbItems = [{ title: 'campaigns', link: '/campaigns' }]

const CampaignsPage = () => {
	const searchParams = useSearchParams()
	const router = useRouter()
	const deleteCampaignMutation = useDeleteCampaignById()
	const updateCampaignByIdMutation = useUpdateCampaignById()

	const page = Number(searchParams.get('page') || 1)
	const pageLimit = Number(searchParams.get('limit') || 0) || 10
	const campaignId = searchParams.get('id')

	const [order] = useState()
	const [status] = useState()
	const [isBusy, setIsBusy] = useState(false)

	const { data: campaignData } = useGetCampaignById(campaignId!, {
		query: {
			enabled: !!campaignId
		}
	})

	const { data: campaignAnalytics } = useGetCampaignAnalyticsById(campaignId!, {
		query: {
			enabled: !!campaignId
		}
	})

	const { data: campaignResponse, refetch: refetchCampaigns } = useGetCampaigns({
		per_page: pageLimit || 10,
		page: page || 1,
		...(order ? { order: order } : {}),
		...(status ? { status: status } : {})
	})

	const totalCampaigns = campaignResponse?.paginationMeta?.total || 0
	const pageCount = Math.ceil(totalCampaigns / pageLimit)
	const campaigns: CampaignSchema[] = campaignResponse?.campaigns || []

	async function deleteCampaign(campaignId: string) {
		try {
			if (!campaignId) return

			setIsBusy(true)

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
		} finally {
			setIsBusy(false)
			await refetchCampaigns()
		}
	}

	async function updateCampaignStatus(
		campaign: CampaignSchema,
		action: 'pause' | 'resume' | 'cancel' | 'running'
	) {
		try {
			setIsBusy(true)
			const response = await updateCampaignByIdMutation.mutateAsync({
				id: campaign.uniqueId,
				data: {
					...campaign,
					status:
						action === 'cancel'
							? 'Cancelled'
							: action === 'pause'
								? 'Paused'
								: 'Running',
					enableLinkTracking: campaign.isLinkTrackingEnabled,
					listIds: campaign.lists.map(list => list.uniqueId),
					tags: campaign.tags.map(tag => tag.uniqueId)
				}
			})

			if (response) {
				// show success notification
				successNotification({
					message: `Campaign ${action === 'cancel' ? 'cancelled' : action} successfully`
				})
			} else {
				// show error notification
				errorNotification({
					message: `Failed to ${action} campaign`
				})
			}
		} catch (error) {
			console.error('Error pausing/resuming campaign', error)
			errorNotification({
				message: 'Error pausing/resuming campaign'
			})
		} finally {
			setIsBusy(false)
			await refetchCampaigns()
		}
	}

	return (
		<>
			<div className="flex-1 space-y-4  p-4 pt-6 md:p-8">
				{campaignId ? (
					<>
						{campaignData && (
							<div className="campaign-details">
								<h2 className="text-2xl font-bold">{campaignData.campaign.name}</h2>
								<p className="text-sm text-gray-600">
									{campaignData.campaign.description}
								</p>
								<div className="mt-4">
									<span className="font-semibold">Start Date:</span>{' '}
									{new Date(campaignData.campaign.createdAt).toLocaleDateString()}
								</div>
								<div>
									<span className="font-semibold">Status:</span>{' '}
									{campaignData.campaign.status}
								</div>
							</div>
						)}

						{/* Campaign Analytics */}
						{campaignAnalytics && (
							<div className="campaign-analytics mt-8">
								<h3 className="text-xl font-bold">Campaign Analytics</h3>
								<div className="mt-4">
									<div>
										<span className="font-semibold">Total Sent:</span>{' '}
										{campaignAnalytics.totalMessages}
									</div>
									<div>
										<span className="font-semibold">Total Delivered:</span>{' '}
										{campaignAnalytics.messagesDelivered}
									</div>
									<div>
										<span className="font-semibold">Total Opened:</span>{' '}
										{campaignAnalytics.messagesRead}
									</div>
									<div>
										<span className="font-semibold">Failed:</span>{' '}
										{campaignAnalytics.messagesFailed}
									</div>
									<div>
										<span className="font-semibold">Un-delivered:</span>{' '}
										{campaignAnalytics.messagesUndelivered}
									</div>
									<div>
										<span className="font-semibold">Sent:</span>{' '}
										{campaignAnalytics.messagesSent}
									</div>
								</div>
							</div>
						)}
					</>
				) : (
					<>
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
							actions={(campaign: CampaignSchema) => {
								const actions: TableCellActionProps[] = []

								// * 1. Edit
								actions.push({
									icon: 'edit',
									label: 'Edit',
									disabled: isBusy,
									onClick: (campaignId: string) => {
										// only allowed if the status is not Scheduled or Draft
										if (
											campaign.status !== 'Scheduled' &&
											campaign.status !== 'Draft'
										) {
											return
										}
										// redirect to the edit page with id in search param
										router.push(`/campaigns/new-or-edit?id=${campaignId}`)
									}
								})

								// * 2. Delete
								actions.push({
									icon: 'trash',
									label: 'Delete',
									disabled: isBusy,
									onClick: () => {
										deleteCampaign(campaign.uniqueId).catch(console.error)
									}
								})

								// * Pause / Resume
								if (campaign.status === 'Running') {
									actions.push({
										icon: 'pause',
										label: 'Pause',
										disabled: isBusy,
										onClick: () => {
											updateCampaignStatus(campaign, 'pause').catch(
												console.error
											)
										}
									})
									// * 3. Cancel
									actions.push({
										icon: 'xCircle',
										label: 'Cancel',
										disabled: isBusy,
										onClick: () => {
											updateCampaignStatus(campaign, 'cancel').catch(
												console.error
											)
										}
									})
								} else if (campaign.status === 'Paused') {
									actions.push({
										icon: 'play',
										label: 'Resume',
										disabled: isBusy,
										onClick: () => {
											updateCampaignStatus(campaign, 'resume').catch(
												console.error
											)
										}
									})
									// * 3. Cancel
									actions.push({
										icon: 'xCircle',
										label: 'Cancel',
										disabled: isBusy,
										onClick: () => {
											updateCampaignStatus(campaign, 'cancel').catch(
												console.error
											)
										}
									})
								} else if (campaign.status === 'Draft') {
									actions.push({
										icon: 'arrowRight',
										label: 'Send',
										disabled: isBusy,
										onClick: () => {
											updateCampaignStatus(campaign, 'running').catch(
												console.error
											)
										}
									})
								}

								return actions
							}}
						/>
					</>
				)}
			</div>
		</>
	)
}

export default CampaignsPage
