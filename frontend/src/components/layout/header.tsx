'use client'

import ThemeToggle from '~/components/layout/ThemeToggle/theme-toggle'
import { clsx as cn } from 'clsx'
import { MobileSidebar } from './mobile-sidebar'
import { UserNav } from './user-nav'
import Link from 'next/link'
import Image from 'next/image'
import { useTheme } from 'next-themes'

export default function Header() {
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
					<UserNav />
					<ThemeToggle />
				</div>
			</nav>
		</div>
	)
}
