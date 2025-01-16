'use client'

import ThemeToggle from '~/components/layout/theme/theme-toggle'
import { clsx as cn } from 'clsx'
import { MobileSidebar } from './mobile-sidebar'
import { UserNav } from './user-nav'
import Link from 'next/link'
import Image from 'next/image'
import { useTheme } from 'next-themes'
import { Notifications } from './notification-dropdown'
import { Button } from '~/components/ui/button'
import { useLayoutStore } from '~/store/layout.store'
import { useAiChatStore } from '~/store/ai-chat-store'

export default function Header() {
	const { featureFlags } = useLayoutStore()
	const { writeProperty } = useAiChatStore()
	const { resolvedTheme } = useTheme()
	return (
		<div className="supports-backdrop-blur:bg-background/60 fixed left-0 right-0 top-0 z-20 border-b bg-background/95 backdrop-blur">
			<nav className="flex h-14 items-center justify-between px-4">
				<div className="hidden lg:block">
					<Link href={'/'}>
						<Image
							src={resolvedTheme === 'dark' ? '/logo/dark.svg' : '/logo/light.svg'}
							height={40}
							width={100}
							alt="logo"
						/>
					</Link>
				</div>
				<div className={cn('block lg:!hidden')}>
					<MobileSidebar />
				</div>

				<div className="flex items-center gap-2">
					{featureFlags?.SystemFeatureFlags.isAiIntegrationEnabled ? (
						<Button
							onClick={() => {
								writeProperty({
									isOpen: true
								})
							}}
						>
							Ask AI
						</Button>
					) : null}

					<Notifications />
					<UserNav />
					<ThemeToggle />
				</div>
			</nav>
		</div>
	)
}
