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
import { ChatBubbleIcon } from '@radix-ui/react-icons'
import { MessageSquareCode, RocketIcon, Phone, InfoIcon } from 'lucide-react'
import { Divider } from '@tremor/react'
import { Toaster } from '~/components/ui/sonner'
import { useGetPrimaryAnalytics, useGetSecondaryAnalytics } from 'root/.generated'
import { type DateRange } from 'react-day-picker'
import React, { useRef, useState } from 'react'
import dayjs from 'dayjs'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '~/components/ui/tooltip'

export default function Page() {
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
									<Divider className="upper text-sm">Campaigns</Divider>
								</CardHeader>
								<CardContent className="flex flex-row items-center justify-between gap-1">
									<div className="flex h-full flex-col gap-2 pt-2">
										<p className="text-xs text-muted-foreground">
											<b>Total</b>:{' '}
											{primaryAnalyticsData?.aggregateAnalytics.campaignStats
												.total || 0}
										</p>
										<p className="text-xs text-muted-foreground">
											<b>Running</b>:{' '}
											{primaryAnalyticsData?.aggregateAnalytics.campaignStats
												.running || 0}
										</p>
									</div>
									<div className="flex h-full flex-col gap-2 pt-2">
										<p className="text-xs text-muted-foreground">
											<b>Draft</b>:{' '}
											{primaryAnalyticsData?.aggregateAnalytics.campaignStats
												.draft || 0}
										</p>
										<p className="text-xs text-muted-foreground">
											<b>Scheduled</b>:{' '}
											{primaryAnalyticsData?.aggregateAnalytics.campaignStats
												.scheduled || 0}
										</p>
									</div>
								</CardContent>
							</Card>
							<Card>
								<CardHeader className="flex flex-row items-center justify-start space-y-0 pb-2">
									<CardTitle className="mx-auto flex w-full flex-row items-center gap-1 text-center text-sm font-medium">
										<ChatBubbleIcon className={`mx-auto size-6`} />
									</CardTitle>
									<Divider className="upper text-sm">Conversations</Divider>
								</CardHeader>
								<CardContent className="flex flex-row items-center justify-between gap-1">
									<div className="flex h-full flex-col gap-2 pt-2">
										<p className="text-xs text-muted-foreground">
											<b>Total</b>:{' '}
											{primaryAnalyticsData?.aggregateAnalytics
												.conversationStats.total || 0}
										</p>
										<p className="text-xs text-muted-foreground">
											<b>Active</b>:{' '}
											{primaryAnalyticsData?.aggregateAnalytics
												.conversationStats.active || 0}
										</p>
									</div>
									<div className="flex h-full flex-col gap-2 pt-2">
										<p className="text-xs text-muted-foreground">
											<b>Resolved</b>:{' '}
											{primaryAnalyticsData?.aggregateAnalytics
												.conversationStats.closed || 0}
										</p>
										<p className="text-xs text-muted-foreground">
											<b>Awaiting Reply</b>:{' '}
											{primaryAnalyticsData?.aggregateAnalytics
												.conversationStats.pending || 0}
										</p>
									</div>
								</CardContent>
							</Card>
							<Card>
								<CardHeader className="flex flex-row items-center justify-start space-y-0 pb-2">
									<CardTitle className="mx-auto flex w-full flex-row items-center gap-1 text-center text-sm font-medium">
										<MessageSquareCode className={`mx-auto size-6`} />
									</CardTitle>
									<Divider className="upper text-sm">Messages</Divider>
								</CardHeader>
								<CardContent className="flex flex-row items-center justify-between gap-1">
									<div className="flex h-full flex-col gap-2 pt-2">
										<p className="text-xs text-muted-foreground">
											<b>Total</b>:{' '}
											{primaryAnalyticsData?.aggregateAnalytics.messageStats
												.total || 0}
										</p>
										<p className="text-xs text-muted-foreground">
											<b>Sent</b>:
											{primaryAnalyticsData?.aggregateAnalytics.messageStats
												.sent || 0}
										</p>
									</div>
									<div className="flex h-full flex-col gap-2 pt-2">
										<p className="text-xs text-muted-foreground">
											<b>Read</b>:
											{primaryAnalyticsData?.aggregateAnalytics.messageStats
												.read || 0}
										</p>
										<p className="text-xs text-muted-foreground">
											<b>Undelivered</b>:
											{primaryAnalyticsData?.aggregateAnalytics.messageStats
												.undelivered || 0}
										</p>
									</div>
								</CardContent>
							</Card>
							<Card>
								<CardHeader className="flex flex-row items-center justify-start space-y-0 pb-2">
									<CardTitle className="mx-auto flex w-full flex-row items-center gap-1 text-center text-sm font-medium">
										<Phone className={`mx-auto size-6`} />
									</CardTitle>
									<Divider className="upper text-sm">Contact</Divider>
								</CardHeader>
								<CardContent className="flex flex-row items-center justify-between gap-1">
									<div className="flex h-full flex-col gap-2 pt-2">
										<p className="text-xs text-muted-foreground">
											<b>Total</b>:
											{primaryAnalyticsData?.aggregateAnalytics.contactStats
												.total || 0}
										</p>
										<p className="text-xs text-muted-foreground">
											<b>Active</b>:
											{primaryAnalyticsData?.aggregateAnalytics.contactStats
												.active || 0}
										</p>
									</div>
									<div className="flex h-full flex-col gap-2 pt-2">
										<p className="text-xs text-muted-foreground">
											<b>Blocked</b>:
											{primaryAnalyticsData?.aggregateAnalytics.contactStats
												.blocked || 0}
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
