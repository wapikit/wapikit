'use client'

import { useRouter } from 'next/navigation'
import { useGetUserNotifications } from 'root/.generated'
import { useAuthState } from '~/hooks/use-auth-state'
import { Icons } from '../icons'
import { Popover, PopoverContent, PopoverTrigger } from '../ui/popover'
import { Separator } from '../ui/separator'

export function Notifications() {
	const router = useRouter()

	const { authState } = useAuthState()

	// ! TODO: use react way point in the dropdown to fetch more notifications
	const { data: notifications } = useGetUserNotifications({
		page: 1,
		per_page: 10
	})
	if (authState.isAuthenticated) {
		return (
			<Popover>
				<PopoverTrigger asChild>
					<Icons.bell className="size-5 text-muted-foreground" />
				</PopoverTrigger>
				<PopoverContent className="w-56 rounded-md" align="end" forceMount>
					<div className="p-3 px-2">Notifications</div>
					<Separator />

					{(notifications?.notifications.length || 0) > 0 ? (
						<div>
							{notifications?.notifications.map(notification => (
								<div
									key={notification.uniqueId}
									className="space-y-2 rounded-md p-2 hover:bg-gray-100"
									onClick={() => {
										if (notification.ctaUrl) {
											router.push(notification.ctaUrl)
										}
									}}
								>
									<p className="text-sm font-medium leading-none">
										{notification.title}
									</p>
									<p className="text-xs leading-none text-muted-foreground">
										{notification.description}
									</p>
									<Separator />
								</div>
							))}
						</div>
					) : (
						<div className="flex min-h-64 items-center justify-center">
							<p className="text-sm font-medium leading-none">No notifications</p>
						</div>
					)}
				</PopoverContent>
			</Popover>
		)
	} else {
		return null
	}
}
