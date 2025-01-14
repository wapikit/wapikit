'use client'

import { LinkClicks } from '~/components/dashboard/link-clicks'
import { CalendarDateRangePicker } from '~/components/date-range-picker'
import { Card, CardContent, CardHeader, CardTitle } from '~/components/ui/card'
import { ScrollArea } from '~/components/ui/scroll-area'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '~/components/ui/tabs'
import { ConversationStatusChart } from '~/components/dashboard/conversation-data'
import { MessageTypeBifurcation } from '~/components/dashboard/message-type-distribution'
import { OrganizationMembers } from '~/components/dashboard/org-members'
import { MessageAggregateAnalytics } from '~/components/dashboard/message-aggregate-stats'
import { ChatBubbleIcon, ExclamationTriangleIcon } from '@radix-ui/react-icons'
import { MessageSquareCode, RocketIcon, Phone, InfoIcon } from 'lucide-react'
import { Callout, Divider } from '@tremor/react'
import { Toaster } from '~/components/ui/sonner'
import { useGetPrimaryAnalytics, useGetSecondaryAnalytics } from 'root/.generated'
import { type DateRange } from 'react-day-picker'
import React, { useRef, useState } from 'react'
import dayjs from 'dayjs'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '~/components/ui/tooltip'
import { useAuthState } from '~/hooks/use-auth-state'
import LoadingSpinner from '~/components/loader'

export default function Page() {
	const { authState } = useAuthState()

	const [date, setDate] = useState<DateRange>({
		from: dayjs().subtract(20, 'day').toDate(),
		to: dayjs().toDate()
	})
	const { data: primaryAnalyticsData } = useGetPrimaryAnalytics({
		from: date.from?.toISOString() || dayjs().subtract(20, 'day').toISOString(),
		to: date.to?.toISOString() || dayjs().toISOString()
	})

	const { data: secondaryAnalyticsData } = useGetSecondaryAnalytics({
		from: date.from?.toISOString() || dayjs().subtract(20, 'day').toISOString(),
		to: date.to?.toISOString() || dayjs().toISOString()
	})

	const datPickerSelectorRef = useRef<HTMLDivElement | null>(null)

	if (authState.isAuthenticated && !authState.data.user.organizationId) {
		return <LoadingSpinner />
	}

	return (
		<ScrollArea className="h-full">
			<Toaster />
			<div className="flex-1 space-y-4 p-4 pt-6 md:p-8">
				<div className="flex items-center justify-between space-y-2">
					<h2 className="text-3xl font-bold tracking-tight">Dashboard</h2>
					<div className="hidden flex-col items-center gap-2 space-x-2 md:flex">
						<TooltipProvider>
							<Tooltip>
								<TooltipTrigger asChild>
									<CalendarDateRangePicker
										dateRange={date}
										setDateRange={setDate}
										ref={datPickerSelectorRef}
									/>
								</TooltipTrigger>
								<TooltipContent
									align="center"
									side="right"
									sideOffset={8}
									className="inline-block"
								>
									<p>
										<InfoIcon /> Select a date range to view analytics data
									</p>
								</TooltipContent>
							</Tooltip>
						</TooltipProvider>
					</div>
				</div>
				<Tabs defaultValue="overview" className="space-y-4">
					<TabsList>
						<TabsTrigger value="overview">Overview</TabsTrigger>
						<TabsTrigger value="conversations">Conversations</TabsTrigger>
					</TabsList>
					<TabsContent value="overview" className="space-y-4">
						<div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
							<Card>
								<CardHeader className="flex flex-row items-center justify-start space-y-0 pb-2">
									<CardTitle className="mx-auto flex w-full flex-row items-center gap-1 text-center text-sm font-medium">
										<RocketIcon className={`mx-auto size-6`} />
									</CardTitle>
									<Divider className="upper text-sm font-bold">Campaigns</Divider>
								</CardHeader>
								<CardContent className="flex flex-row items-center justify-between gap-1">
									<div className="flex h-full flex-col gap-2 pt-2">
										<p className="text-sm font-light text-muted-foreground">
											<b>Total</b>:{' '}
											<span className="font-extrabold">
												{primaryAnalyticsData?.aggregateAnalytics
													.campaignStats.totalCampaigns || 0}
											</span>
										</p>
										<p className="text-sm font-light text-muted-foreground">
											<b>Running</b>:{' '}
											<span className="font-extrabold">
												{primaryAnalyticsData?.aggregateAnalytics
													.campaignStats.campaignsRunning || 0}
											</span>
										</p>
									</div>
									<div className="flex h-full flex-col gap-2 pt-2">
										<p className="text-sm font-light text-muted-foreground">
											<b>Draft</b>:{' '}
											<span className="font-extrabold">
												{primaryAnalyticsData?.aggregateAnalytics
													.campaignStats.campaignsDraft || 0}
											</span>
										</p>
										<p className="text-sm font-light text-muted-foreground">
											<b>Scheduled</b>:{' '}
											<span className="font-extrabold">
												{primaryAnalyticsData?.aggregateAnalytics
													.campaignStats.campaignsScheduled || 0}
											</span>
										</p>
									</div>
								</CardContent>
							</Card>
							<Card>
								<CardHeader className="flex flex-row items-center justify-start space-y-0 pb-2">
									<CardTitle className="mx-auto flex w-full flex-row items-center gap-1 text-center text-sm font-medium">
										<ChatBubbleIcon className={`mx-auto size-6`} />
									</CardTitle>
									<Divider className="upper text-sm font-bold">
										Conversations
									</Divider>
								</CardHeader>
								<CardContent className="flex flex-row items-center justify-between gap-1">
									<div className="flex h-full flex-col gap-2 pt-2">
										<p className="text-sm font-light text-muted-foreground">
											<b>Total</b>:{' '}
											<span className="font-extrabold">
												{primaryAnalyticsData?.aggregateAnalytics
													.conversationStats.totalConversations || 0}
											</span>
										</p>
										<p className="text-sm font-light text-muted-foreground">
											<b>Active</b>:{' '}
											<span className="font-extrabold">
												{primaryAnalyticsData?.aggregateAnalytics
													.conversationStats.conversationsActive || 0}
											</span>
										</p>
									</div>
									<div className="flex h-full flex-col gap-2 pt-2">
										<p className="text-sm font-light text-muted-foreground">
											<b>Resolved</b>:{' '}
											<span className="font-extrabold">
												{primaryAnalyticsData?.aggregateAnalytics
													.conversationStats.conversationsClosed || 0}
											</span>
										</p>
										<p className="text-sm font-light text-muted-foreground">
											<b>Awaiting Reply</b>:{' '}
											<span className="font-extrabold">
												{primaryAnalyticsData?.aggregateAnalytics
													.conversationStats.conversationsPending || 0}
											</span>
										</p>
									</div>
								</CardContent>
							</Card>
							<Card>
								<CardHeader className="flex flex-row items-center justify-start space-y-0 pb-2">
									<CardTitle className="mx-auto flex w-full flex-row items-center gap-1 text-center text-sm font-medium">
										<MessageSquareCode className={`mx-auto size-6`} />
									</CardTitle>
									<Divider className="upper text-sm font-bold">Messages</Divider>
								</CardHeader>
								<CardContent className="flex flex-row items-center justify-between gap-1">
									<div className="flex h-full flex-col gap-2 pt-2">
										<p className="text-sm font-light text-muted-foreground">
											<b>Total</b>:{' '}
											<span className="font-extrabold">
												{primaryAnalyticsData?.aggregateAnalytics
													.messageStats.totalMessages || 0}
											</span>
										</p>
										<p className="text-sm font-light text-muted-foreground">
											<b>Sent</b>:
											<span className="font-extrabold">
												{primaryAnalyticsData?.aggregateAnalytics
													.messageStats.messagesSent || 0}
											</span>
										</p>
									</div>
									<div className="flex h-full flex-col gap-2 pt-2">
										<p className="text-sm font-light text-muted-foreground">
											<b>Read</b>:
											<span className="font-extrabold">
												{primaryAnalyticsData?.aggregateAnalytics
													.messageStats.messagesRead || 0}
											</span>
										</p>
										<p className="text-sm font-light text-muted-foreground">
											<b>Undelivered</b>:
											<span className="font-extrabold">
												{primaryAnalyticsData?.aggregateAnalytics
													.messageStats.messagesUndelivered || 0}
											</span>
										</p>
									</div>
								</CardContent>
							</Card>
							<Card>
								<CardHeader className="flex flex-row items-center justify-start space-y-0 pb-2">
									<CardTitle className="mx-auto flex w-full flex-row items-center gap-1 text-center text-sm font-medium">
										<Phone className={`mx-auto size-6`} />
									</CardTitle>
									<Divider className="upper text-sm font-bold">Contacts</Divider>
								</CardHeader>
								<CardContent className="flex flex-row items-center justify-between gap-1">
									<div className="flex h-full flex-col gap-2 pt-2">
										<p className="text-sm font-light text-muted-foreground">
											<b>Total</b>:
											<span className="font-extrabold">
												{primaryAnalyticsData?.aggregateAnalytics
													.contactStats.totalContacts || 0}
											</span>
										</p>
										<p className="text-sm font-light text-muted-foreground">
											<b>Active</b>:
											<span className="font-extrabold">
												{primaryAnalyticsData?.aggregateAnalytics
													.contactStats.contactsActive || 0}
											</span>
										</p>
									</div>
									<div className="flex h-full flex-col gap-2 pt-2">
										<p className="text-sm font-light text-muted-foreground">
											<b>Blocked</b>:
											<span className="font-extrabold">
												{primaryAnalyticsData?.aggregateAnalytics
													.contactStats.contactsBlocked || 0}
											</span>
										</p>
									</div>
								</CardContent>
							</Card>
						</div>
						<div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-8">
							<Card className="col-span-2 md:col-span-4">
								<CardHeader>
									<CardTitle>Message Analytics</CardTitle>
								</CardHeader>
								<CardContent className="pl-2">
									<MessageAggregateAnalytics
										data={primaryAnalyticsData?.messageAnalytics || []}
									/>
								</CardContent>
							</Card>
							<Card className="col-span-4 md:col-span-4">
								<CardHeader>
									<CardTitle>Link Clicks</CardTitle>
								</CardHeader>
								<CardContent className="pl-2">
									<LinkClicks
										data={primaryAnalyticsData?.linkClickAnalytics || []}
									/>
								</CardContent>
							</Card>
						</div>
					</TabsContent>
					<TabsContent value="conversations" className="space-y-4">
						<Callout title="" icon={ExclamationTriangleIcon}>
							{' '}
							These analytics will be available, once the{' '}
							<a href="/conversations" className="underline">
								live team inbox conversation
							</a>{' '}
							feature will be shipped soon.{' '}
						</Callout>
						<div>
							<Card className="col-span-2 md:col-span-4">
								<CardHeader>
									<CardTitle>Message Type Distribution</CardTitle>
								</CardHeader>
								<CardContent className="pl-2">
									<MessageTypeBifurcation
										data={
											secondaryAnalyticsData?.messageTypeTrafficDistributionAnalytics ||
											[]
										}
									/>
								</CardContent>
							</Card>
						</div>
						<div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-8">
							<Card className="col-span-3 md:col-span-4">
								<CardHeader>
									<CardTitle>Organization Members</CardTitle>
								</CardHeader>
								<CardContent className="pl-2">
									<OrganizationMembers />
								</CardContent>
							</Card>
							<Card className="col-span-3 md:col-span-4">
								<CardHeader>
									<CardTitle>Conversation Status</CardTitle>
								</CardHeader>
								<CardContent className="pl-2">
									<ConversationStatusChart />
								</CardContent>
							</Card>
						</div>
					</TabsContent>
				</Tabs>
			</div>
		</ScrollArea>
	)
}
