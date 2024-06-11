import { LinkClicks } from '~/components/dashboard/link-clicks'
import { CalendarDateRangePicker } from '~/components/date-range-picker'
import { Button } from '~/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '~/components/ui/card'
import { ScrollArea } from '~/components/ui/scroll-area'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '~/components/ui/tabs'
import { ConversationStatusChart } from '~/components/dashboard/conversation-data'
import { MessageTypeBifurcation } from '~/components/dashboard/message-type-distribution'
import { OrganizationMembers } from '~/components/dashboard/org-members'
import { MessageAggregateAnalytics } from '~/components/dashboard/message-aggregate-stats'
import { ChatBubbleIcon } from '@radix-ui/react-icons'
import { MessageSquareCode, RocketIcon, Phone } from 'lucide-react'
import { Divider } from '@tremor/react'

export default function page() {
	return (
		<ScrollArea className="h-full ">
			<div className="flex-1 space-y-4 p-4 pt-6 md:p-8">
				<div className="flex items-center justify-between space-y-2">
					<h2 className="text-3xl font-bold tracking-tight">Dashboard</h2>
					<div className="hidden items-center space-x-2 md:flex">
						<CalendarDateRangePicker />
						<Button>View</Button>
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
											<b>Total</b>: 0
										</p>
										<p className="text-xs text-muted-foreground">
											<b>Running</b>: 0
										</p>
									</div>
									<div className="flex h-full flex-col gap-2 pt-2">
										<p className="text-xs text-muted-foreground">
											<b>Draft</b>: 0
										</p>
										<p className="text-xs text-muted-foreground">
											<b>Scheduled</b>: 0
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
											<b>Total</b>: 0
										</p>
										<p className="text-xs text-muted-foreground">
											<b>Open</b>: 0
										</p>
									</div>
									<div className="flex h-full flex-col gap-2 pt-2">
										<p className="text-xs text-muted-foreground">
											<b>Resolved</b>: 0
										</p>
										<p className="text-xs text-muted-foreground">
											<b>Awaiting Reply</b>: 0
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
											<b>Total</b>: 0
										</p>
										<p className="text-xs text-muted-foreground">
											<b>Read</b>: 0
										</p>
									</div>
									<div className="flex h-full flex-col gap-2 pt-2">
										<p className="text-xs text-muted-foreground">
											<b>Delivered</b>: 0
										</p>
										<p className="text-xs text-muted-foreground">
											<b>Scheduled</b>: 0
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
											<b>Total</b>: 0
										</p>
										<p className="text-xs text-muted-foreground">
											<b>Running</b>: 0
										</p>
									</div>
									<div className="flex h-full flex-col gap-2 pt-2">
										<p className="text-xs text-muted-foreground">
											<b>Draft</b>: 0
										</p>
										<p className="text-xs text-muted-foreground">
											<b>Scheduled</b>: 0
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
									<MessageAggregateAnalytics />
								</CardContent>
							</Card>
							<Card className="col-span-4 md:col-span-4">
								<CardHeader>
									<CardTitle>Link Clicks</CardTitle>
								</CardHeader>
								<CardContent className="pl-2">
									<LinkClicks />
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
									<MessageTypeBifurcation />
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
