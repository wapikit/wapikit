import { CalendarDateRangePicker } from '~/components/date-range-picker'
import { Overview } from '~/components/overview'
import { RecentSales } from '~/components/recent-sales'
import { Button } from '~/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '~/components/ui/card'
import { ScrollArea } from '~/components/ui/scroll-area'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '~/components/ui/tabs'

export default function page() {
	return (
		<ScrollArea className="h-full">
			<div className="flex-1 space-y-4 p-4 pt-6 md:p-8">
				<div className="flex items-center justify-between space-y-2">
					<h2 className="text-3xl font-bold tracking-tight">Dashboard</h2>
					<div className="hidden items-center space-x-2 md:flex">
						<CalendarDateRangePicker />
						<Button>Download</Button>
					</div>
				</div>
				<Tabs defaultValue="overview" className="space-y-4">
					<TabsList>
						<TabsTrigger value="overview">Overview</TabsTrigger>
						<TabsTrigger value="analytics">Analytics</TabsTrigger>
					</TabsList>
					<TabsContent value="overview" className="space-y-4">
						<div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
							<Card>
								<CardHeader className="flex flex-row items-center justify-start space-y-0 pb-2">
									<CardTitle className="flex w-full flex-row gap-1 text-sm font-medium">
										<svg
											xmlns="http://www.w3.org/2000/svg"
											viewBox="0 0 24 24"
											fill="none"
											stroke="currentColor"
											strokeLinecap="round"
											strokeLinejoin="round"
											strokeWidth="2"
											className="h-4 w-4 text-muted-foreground"
										>
											<path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6" />
										</svg>
										Campaigns
									</CardTitle>
								</CardHeader>
								<CardContent className="flex flex-row items-center justify-between gap-1 space-y-2">
									<div>
										<p className="text-xs text-muted-foreground">
											<b>Total</b>: 0
										</p>
										<p className="text-xs text-muted-foreground">
											<b>Running</b>: 0
										</p>
									</div>
									<div>
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
									<CardTitle className="flex w-full flex-row text-sm font-medium">
										<svg
											xmlns="http://www.w3.org/2000/svg"
											viewBox="0 0 24 24"
											fill="none"
											stroke="currentColor"
											strokeLinecap="round"
											strokeLinejoin="round"
											strokeWidth="2"
											className="h-4 w-4 text-muted-foreground"
										>
											<path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6" />
										</svg>
										Conversations
									</CardTitle>
								</CardHeader>
								<CardContent className="flex flex-row items-center justify-between gap-1 space-y-2">
									<div>
										<p className="text-xs text-muted-foreground">
											<b>Total</b>: 0
										</p>
										<p className="text-xs text-muted-foreground">
											<b>Open</b>: 0
										</p>
									</div>
									<div>
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
								<CardHeader className="flex flex-row items-center justify-start gap-1 space-y-0 pb-2">
									<CardTitle className="flex w-full flex-row gap-1 text-sm font-medium">
										<svg
											xmlns="http://www.w3.org/2000/svg"
											viewBox="0 0 24 24"
											fill="none"
											stroke="currentColor"
											strokeLinecap="round"
											strokeLinejoin="round"
											strokeWidth="2"
											className="h-4 w-4 text-muted-foreground"
										>
											<path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6" />
										</svg>
										Messages
									</CardTitle>
								</CardHeader>
								<CardContent className="flex flex-row items-center justify-between gap-1 space-y-2">
									<div>
										<p className="text-xs text-muted-foreground">
											<b>Total</b>: 0
										</p>
										<p className="text-xs text-muted-foreground">
											<b>Read</b>: 0
										</p>
									</div>
									<div>
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
									<CardTitle className="flex w-full flex-row text-sm font-medium">
										<svg
											xmlns="http://www.w3.org/2000/svg"
											viewBox="0 0 24 24"
											fill="none"
											stroke="currentColor"
											strokeLinecap="round"
											strokeLinejoin="round"
											strokeWidth="2"
											className="h-4 w-4 text-muted-foreground"
										>
											<path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6" />
										</svg>
										Campaigns
									</CardTitle>
								</CardHeader>
								<CardContent className="flex flex-row items-center justify-between gap-1 space-y-2">
									<div>
										<p className="text-xs text-muted-foreground">
											<b>Total</b>: 0
										</p>
										<p className="text-xs text-muted-foreground">
											<b>Running</b>: 0
										</p>
									</div>
									<div>
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
						<div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-7">
							<Card className="col-span-4">
								<CardHeader>
									<CardTitle>Overview</CardTitle>
								</CardHeader>
								<CardContent className="pl-2">
									<Overview />
								</CardContent>
							</Card>
							<Card className="col-span-4 md:col-span-3">
								<CardHeader>
									<CardTitle>Organizartion Members</CardTitle>
									<CardDescription>6 members online</CardDescription>
								</CardHeader>
								<CardContent>
									<RecentSales />
								</CardContent>
							</Card>
						</div>
					</TabsContent>
				</Tabs>
			</div>
		</ScrollArea>
	)
}
