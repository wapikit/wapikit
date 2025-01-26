'use client'

import React, { useCallback, useEffect, useRef, useState } from 'react'
import { DashboardNav } from '~/components/dashboard-nav'
import { navItems, pathAtWhichSidebarShouldBeCollapsedByDefault } from '~/constants'
import { clsx as cn } from 'clsx'
import { ChevronLeft } from 'lucide-react'
import { useSidebar } from '~/hooks/use-sidebar'
import { usePathname } from 'next/navigation'

type SidebarProps = {
	className?: string
}

export default function Sidebar({ className }: SidebarProps) {
	const { isMinimized, toggle } = useSidebar()
	const [status, setStatus] = useState(false)

	const isNavbarCollapsedOnce = useRef(false)

	const pathname = usePathname()

	const handleToggle = useCallback(() => {
		setStatus(true)
		toggle()
		setTimeout(() => setStatus(false), 500), []
	}, [toggle])

	// ! COMMAND + B to toggle sidebar
	useEffect(() => {
		const handleKeydown = (event: KeyboardEvent) => {
			const isMac = navigator.platform.toUpperCase().includes('MAC')
			const isCommandOrCtrlPressed = isMac ? event.metaKey : event.ctrlKey

			if (isCommandOrCtrlPressed && event.key === 'b') {
				event.preventDefault() // Prevent default browser behavior
				handleToggle()
			}
		}

		window.addEventListener('keydown', handleKeydown)

		return () => {
			window.removeEventListener('keydown', handleKeydown)
		}
	}, [handleToggle])

	useEffect(() => {
		if (isNavbarCollapsedOnce.current) return
		const shouldBeCollapsed = pathAtWhichSidebarShouldBeCollapsedByDefault.some(path =>
			pathname.includes(path)
		)

		if (!isMinimized && shouldBeCollapsed) {
			handleToggle()
			isNavbarCollapsedOnce.current = true
		}
	}, [handleToggle, isMinimized, pathname, toggle])

	return (
		<nav
			className={cn(
				`relative hidden h-screen flex-none border-r pt-12 md:block`,
				status && 'duration-500',
				!isMinimized ? 'w-72' : 'w-[72px]',
				className
			)}
		>
			<ChevronLeft
				className={cn(
					'absolute -right-3 top-16 cursor-pointer rounded-full border bg-background text-3xl text-foreground',
					isMinimized && 'rotate-180'
				)}
				onClick={handleToggle}
			/>
			<div className="space-y-4 py-4">
				<div className="px-3 py-2">
					<div className="mt-3 space-y-1">
						<DashboardNav items={navItems} />
					</div>
				</div>
			</div>
		</nav>
	)
}
