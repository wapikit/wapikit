'use client'

import BreadCrumb from '~/components/breadcrumb'
import { CampaignTableColumns } from '~/components/tables/columns'
import { TableComponent } from '~/components/tables/table'
import { Button, buttonVariants } from '~/components/ui/button'
import { Heading } from '~/components/ui/heading'
import { Separator } from '~/components/ui/separator'
import {
	CampaignStatusEnum,
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
import { useRouter, useSearchParams } from 'next/navigation'
import { useState } from 'react'
import { errorNotification, materialConfirm, successNotification } from '~/reusable-functions'
import type { TableCellActionProps } from '~/types'
import { Card, CardContent, CardHeader, CardTitle } from '~/components/ui/card'
import { Badge } from '~/components/ui/badge'
import dayjs from 'dayjs'
import { Icons } from '~/components/icons'
import { LinkClicks } from '~/components/dashboard/link-clicks'
import { MessageAggregateAnalytics } from '~/components/dashboard/message-aggregate-stats'
import { ScrollArea } from '~/components/ui/scroll-area'

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

	// eslint-disable-next-line @typescript-eslint/no-non-null-assertion
	const { data: campaignData, refetch: refetchCampaign } = useGetCampaignById(campaignId!, {
		query: {
			enabled: !!campaignId
		}
	})

	// eslint-disable-next-line @typescript-eslint/no-non-null-assertion
	const { data: campaignAnalytics } = useGetCampaignAnalyticsById(campaignId!, {
		query: {
			enabled: !!campaignId,
			...(campaignData && campaignData?.campaign.status === CampaignStatusEnum.Running
				? {
						refetchInterval: 10000
					}
				: {})
		}
	})

	const { data: campaignResponse, refetch: refetchCampaigns } = useGetCampaigns(
		{
			per_page: pageLimit || 10,
			page: page || 1,
			...(order ? { order: order } : {}),
			...(status ? { status: status } : {})
		},
		{
			query: {
				enabled: !campaignId
			}
		}
	)

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

				if (campaignId) {
					await refetchCampaign()
				}
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

			const confirmation = await materialConfirm({
				title: `${action === 'cancel' ? 'Cancel' : action === 'pause' ? 'Pause' : action === 'resume' ? 'Resume' : 'Send'} Campaign`,
				description: `Are you sure you want to ${
					action === 'running' ? 'start' : action
				} this campaign?`
			})

			if (!confirmation) return

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
				if (campaignId) {
					await refetchCampaign()
				}
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
		<ScrollArea className="h-full">
			<div className="flex-1 space-y-4  p-4 pt-6 md:p-8">
				{campaignId ? (
					<>
						<BreadCrumb items={breadcrumbItems} />
						<div className="flex items-center justify-between">
							<Heading title={`Campaigns Details`} description="" />
							{campaignData && (
								<div className="flex w-fit flex-row items-center justify-end gap-3 p-6">
									<Button
										onClick={() => {
											router.push(
												`/campaigns/new-or-edit?id=${campaignData.campaign.uniqueId}`
											)
										}}
										disabled={
											isBusy ||
											campaignData.campaign.status === 'Running' ||
											campaignData.campaign.status === 'Cancelled' ||
											campaignData.campaign.status === 'Finished'
										}
										className="flex flex-row gap-2"
									>
										<Icons.edit className="size-4" />
										Edit
									</Button>
									<Button
										variant={'destructive'}
										disabled={
											isBusy || campaignData.campaign.status === 'Running'
										}
										onClick={() => {
											deleteCampaign(campaignData.campaign.uniqueId).catch(
												console.error
											)
										}}
										className="flex flex-row gap-2"
									>
										<Icons.edit className="size-4" />
										Delete
									</Button>
									{campaignData.campaign.status === 'Running' ? (
										<>
											<Button
												variant={'secondary'}
												onClick={() => {
													updateCampaignStatus(
														campaignData.campaign,
														'pause'
													).catch(console.error)
												}}
												className="flex flex-row gap-2"
											>
												<Icons.pause className="size-4" />
												Pause
											</Button>

											<Button
												variant={'secondary'}
												onClick={() => {
													updateCampaignStatus(
														campaignData.campaign,
														'cancel'
													).catch(console.error)
												}}
												className="flex flex-row gap-2"
											>
												<Icons.xCircle className="size-4" />
												Cancel
											</Button>
										</>
									) : campaignData.campaign.status === 'Paused' ? (
										<>
											<Button
												onClick={() => {
													updateCampaignStatus(
														campaignData.campaign,
														'resume'
													).catch(console.error)
												}}
												className="flex flex-row gap-2"
											>
												<Icons.play className="size-4" />
												Resume
											</Button>

											<Button
												variant={'secondary'}
												onClick={() => {
													updateCampaignStatus(
														campaignData.campaign,
														'cancel'
													).catch(console.error)
												}}
												className="flex flex-row gap-2"
											>
												<Icons.xCircle className="size-4" />
												Cancel
											</Button>
										</>
									) : campaignData.campaign.status === 'Draft' ? (
										<>
											<Button
												onClick={() => {
													updateCampaignStatus(
														campaignData.campaign,
														'running'
													).catch(console.error)
												}}
												className="flex flex-row gap-2"
											>
												<Icons.arrowRight className="size-4" />
												Send
											</Button>
										</>
									) : null}
								</div>
							)}
						</div>
						{campaignData && (
							<Card className="flex flex-col gap-4">
								<CardContent className="mt-4 flex flex-row items-start justify-between gap-2">
									<div className="flex h-full flex-1 flex-col items-start justify-start gap-2 rounded-md border p-4">
										<div className="flex flex-row items-center">
											<span className="text-sm font-semibold">
												name :&nbsp;{' '}
											</span>{' '}
											<div className="flex flex-wrap items-center justify-center gap-4">
												{campaignData.campaign.name}
												<Badge
													variant={
														campaignData.campaign.status === 'Draft'
															? 'outline'
															: campaignData.campaign.status ===
																  'Cancelled'
																? 'destructive'
																: 'default'
													}
													className={clsx(
														campaignData.campaign.status === 'Paused' ||
															campaignData.campaign.status ===
																'Scheduled'
															? 'bg-yellow-500'
															: campaignData.campaign.status ===
																  'Cancelled'
																? 'bg-red-300'
																: ''
													)}
												>
													{campaignData.campaign.status}
												</Badge>
												{campaignData.campaign.status === 'Running' ? (
													<div className="flex h-full w-fit items-center justify-center">
														<div className="rotate h-4 w-4 animate-spin rounded-full border-4 border-solid  border-l-primary" />
													</div>
												) : null}
											</div>
										</div>

										<p className="text-sm text-muted-foreground">
											description: {campaignData.campaign.description}
										</p>
										<div className="flex h-full flex-col gap-6 pt-2">
											{/* sent to lists */}
											<div className="flex flex-row items-center">
												<span className="text-sm font-semibold">
													Sent To:&nbsp;{' '}
												</span>{' '}
												<div className="flex flex-wrap items-center justify-center gap-0.5 truncate">
													{campaignData.campaign.lists.length === 0 && (
														<Badge variant={'outline'}>None</Badge>
													)}
													{campaignData.campaign.lists.map(list => {
														return (
															<Badge key={list.uniqueId}>
																{list.name}
															</Badge>
														)
													})}
												</div>
											</div>

											{/* tags */}
											<div className="flex flex-row items-center">
												<span className="text-sm font-semibold">
													Tags:&nbsp;{' '}
												</span>{' '}
												<div className="flex flex-wrap items-center justify-center gap-0.5 truncate">
													{campaignData.campaign.tags.length === 0 && (
														<Badge variant={'outline'}>None</Badge>
													)}
													{campaignData.campaign.tags.map(tag => {
														return (
															<Badge key={tag.uniqueId}>
																{tag.name}
															</Badge>
														)
													})}
												</div>
											</div>

											{/* created on */}
											<div className="flex flex-row items-center">
												<span className="text-sm font-semibold">
													created on:&nbsp;{' '}
												</span>{' '}
												<div className="flex flex-wrap items-center justify-center gap-0.5 truncate">
													{dayjs(campaignData.campaign.createdAt).format(
														'DD MMM, YYYY'
													)}
												</div>
											</div>
										</div>
									</div>
									<div className="flex h-full flex-1 flex-col rounded-md border p-4">
										<div className="flex h-full w-full flex-col gap-2 pt-2">
											<p className="flex flex-row text-sm font-light text-muted-foreground">
												<span className="flex gap-2">
													<Icons.check className="size-5" />
													<b>Messages Sent:</b>{' '}
												</span>
												<span className="font-extrabold">
													{campaignAnalytics?.totalMessages || 0}
												</span>
											</p>
											<p className="flex flex-row text-sm font-light text-muted-foreground">
												<span className="flex gap-2">
													<Icons.doubleCheck className="size-5" />
													<b>Messages Delivered:</b>{' '}
												</span>
												<span className="font-extrabold">
													{campaignAnalytics?.messagesDelivered || 0}
												</span>
											</p>
											<p className="flex flex-row text-sm font-light text-muted-foreground">
												<span className="flex gap-2">
													<Icons.doubleCheck className="size-5 text-primary" />
													<b>Messages Read:</b>{' '}
												</span>
												<span className="font-extrabold">
													{campaignAnalytics?.messagesRead || 0}
												</span>
											</p>
											<p className="flex flex-row text-sm font-light text-muted-foreground">
												<span className="flex gap-2">
													<Icons.xCircle className="size-5" />
													<b>Messages Failed:</b>{' '}
												</span>
												<span className="font-extrabold">
													{campaignAnalytics?.messagesFailed || 0}
												</span>
											</p>
											<p className="flex flex-row text-sm font-light text-muted-foreground">
												<span className="flex gap-2">
													<Icons.xCircle className="size-5" />
													<b>Messages Undelivered:</b>{' '}
												</span>
												<span className="font-extrabold">
													{campaignAnalytics?.messagesUndelivered || 0}
												</span>
											</p>
										</div>
									</div>
								</CardContent>
							</Card>
						)}
						{campaignAnalytics && (
							<div className="flex flex-row gap-4 ">
								<Card className="flex-1">
									<CardHeader>
										<CardTitle>Message Analytics</CardTitle>
									</CardHeader>
									<CardContent className="pl-2">
										<MessageAggregateAnalytics data={[]} />
									</CardContent>
								</Card>
								{campaignData?.campaign.isLinkTrackingEnabled && (
									<Card className="flex-1">
										<CardHeader>
											<CardTitle>Link Clicks</CardTitle>
										</CardHeader>
										<CardContent className="pl-2">
											<LinkClicks data={[]} />
										</CardContent>
									</Card>
								)}
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
		</ScrollArea>
	)
}

export default CampaignsPage
