'use client'

import { Command } from 'cmdk'
import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { useLayoutStore } from '~/store/layout.store'
import { motion } from 'framer-motion'
import {
	Dialog,
	DialogContent,
	DialogTrigger,
	DialogPortal,
	DialogOverlay
} from '~/components/ui/dialog'
import { VisuallyHidden } from '@radix-ui/react-visually-hidden'
import { Icons } from '../icons'
import { Badge } from '../ui/badge'
import { useRouter } from 'next/navigation'
import { type CommandItemType } from '~/types'
import { clsx } from 'clsx'

export default function CommandMenuProvider() {
	const router = useRouter()

	const [input, setInput] = useState('')

	const commandItemsAndGroups: {
		groupLabel: '' | 'Campaigns' | 'Contacts' | 'Teams'
		items: CommandItemType[]
	}[] = useMemo(() => {
		return [
			{
				groupLabel: '',
				items: [
					{
						icon: 'sparkles',
						label: 'Ask AI',
						action: () => {
							router.push(`/ai?s=true`)
						},
						slug: 'ask-ai'
					},
					{
						icon: 'settings',
						label: 'Settings',
						action: () => {
							router.push('/settings')
						},
						slug: 'settings'
					},
					{
						icon: 'page',
						label: 'Documentation',
						action: () => {
							router.push('/docs')
						},
						slug: 'docs'
					}
				]
			},
			{
				groupLabel: 'Campaigns',
				items: [
					{
						icon: 'add',
						label: 'Create Campaign',
						action: () => {
							router.push('/campaigns/new-or-edit')
						},
						slug: 'create-campaign'
					}
				]
			},
			{
				groupLabel: 'Contacts',
				items: [
					{
						icon: 'add',
						label: 'Create List',
						action: () => {
							router.push('/lists/new-or-edit')
						},
						slug: 'create-list'
					},
					{
						icon: 'download',
						label: 'Bulk Import Contacts',
						action: () => {
							router.push('/contacts/bulk-import')
						},
						slug: 'bulk-import-contacts'
					},
					{
						icon: 'add',
						label: 'Create Contact',
						action: () => {
							router.push('/contacts/new-or-edit')
						},
						slug: 'create-contact'
					}
				]
			},
			{
				groupLabel: 'Teams',
				items: [
					{
						icon: 'user',
						label: 'Invite team member',
						action: () => {
							router.push('/team')
						},
						slug: 'invite-team-member'
					}
				]
			}
		]
	}, [router])

	const { isCommandMenuOpen, writeProperty } = useLayoutStore()

	const scrollContainerRef = useRef<HTMLDivElement>(null)

	const [currentSelected, setCurrentSelected] = useState<string>('ask-ai')

	const runAction = useCallback(
		(action: () => void, slug: string) => {
			writeProperty({
				isCommandMenuOpen: false
			})
			setCurrentSelected(() => slug)
			action()
		},
		[writeProperty]
	)

	useEffect(() => {
		const down = (e: KeyboardEvent) => {
			if (e.key === 'k' && (e.metaKey || e.ctrlKey)) {
				e.preventDefault()
				writeProperty({
					isCommandMenuOpen: !isCommandMenuOpen
				})
			}

			if (e.key === 'Escape') {
				writeProperty({
					isCommandMenuOpen: false
				})
			}

			// key up and down handle
			if (e.key === 'ArrowDown') {
				setCurrentSelected((currentSlug: string) => {
					const allItems = commandItemsAndGroups
						.map(group => group.items.map(item => item.slug))
						.flat()

					const currentIndex = allItems.indexOf(currentSlug)

					if (currentIndex === -1) {
						return allItems[0]
					} else {
						if (currentIndex + 1 < allItems.length) {
							return allItems[currentIndex + 1]
						} else {
							return allItems[0]
						}
					}
				})
			}

			if (e.key === 'ArrowUp') {
				setCurrentSelected(currentSlug => {
					const allItems = commandItemsAndGroups
						.map(group => group.items.map(item => item.slug))
						.flat()

					const currentIndex = allItems.indexOf(currentSlug)

					if (currentIndex === -1) {
						return allItems[0]
					} else {
						if (currentIndex - 1 >= 0) {
							return allItems[currentIndex - 1]
						} else {
							return allItems[allItems.length - 1]
						}
					}
				})
			}

			if (e.key === 'Enter') {
				const item = commandItemsAndGroups
					.map(group => group.items)
					.flat()
					.find(item => item.slug === currentSelected)

				if (item) {
					runAction(item.action, item.slug)
				} else {
					if (input) {
						router.push(`/ai?question=${encodeURIComponent(input)}`)
					} else {
						// IMPOSSIBLE CASE
					}
				}
			}
		}

		document.addEventListener('keydown', down)
		return () => document.removeEventListener('keydown', down)
	}, [
		writeProperty,
		isCommandMenuOpen,
		commandItemsAndGroups,
		currentSelected,
		runAction,
		router,
		input
	])

	useEffect(() => {
		const selectedItem = document.querySelector(
			`[cmkd-item="${currentSelected}"]`
		) as unknown as HTMLElement | null

		const container = scrollContainerRef.current

		if (selectedItem && container) {
			const itemTop = selectedItem.offsetTop
			const itemHeight = selectedItem.offsetHeight
			const containerTop = container.scrollTop
			const containerHeight = container.offsetHeight

			// Calculate scroll position
			if (itemTop < containerTop) {
				// Item is above visible area
				container.scrollTo({
					top: itemTop,
					behavior: 'smooth'
				})

				console.log('scrolling up')
			} else if (itemTop + itemHeight > containerTop + containerHeight) {
				// Item is below visible area
				container.scrollTo({
					top: itemTop - containerHeight + itemHeight,
					behavior: 'smooth'
				})

				console.log('scrolling down')
			}
		}
	}, [currentSelected])

	return (
		<Dialog
			onOpenChange={isOpen => {
				writeProperty({
					isCommandMenuOpen: isOpen
				})
			}}
			open={isCommandMenuOpen}
			key={'command_menu'}
		>
			<VisuallyHidden>
				<DialogTrigger />
			</VisuallyHidden>

			<DialogPortal>
				<DialogOverlay className="opacity-5" />
				<DialogContent className="!p-0">
					<motion.div
						initial={{ opacity: 0, scale: 0.98 }}
						animate={{ opacity: 1, scale: 1 }}
						exit={{ opacity: 0, scale: 0.98 }}
						transition={{ duration: 0.2 }}
					>
						<div className="command-menu relative z-[200]">
							<Command className="flex max-h-96 flex-col gap-2">
								<div className="relative h-fit">
									<Command.Input
										autoFocus
										placeholder="Search or Ask Anything..."
										className="focus:ring-accent"
										value={input}
										onValueChange={value => {
											setInput(value)
										}}
									/>

									<div className="absolute right-4 top-2 flex gap-2">
										<div className="flex gap-1">
											<kbd className="rounded-md bg-accent px-2 py-1 text-xs">
												Esc
											</kbd>
										</div>
									</div>
								</div>
								<Command.Empty className="mx-2 flex !h-auto items-center justify-between gap-3 rounded-md bg-accent px-3 py-3">
									<div className="items-star mb-auto flex h-fit flex-1 flex-row justify-start gap-3">
										<Badge className="flex w-20 flex-row gap-1 px-2 py-1 text-xs">
											Ask AI
											<Icons.sparkles size={12} />
										</Badge>
										<p className="max-w-sm flex-1 whitespace-normal break-words">
											{input}
										</p>
									</div>
									<div className="opacity-100 transition-all ease-in-out">
										<span className="inline-flex h-5 w-fit items-center justify-center rounded-md px-2 py-0.5 text-center text-xs font-semibold leading-4 text-gray-700 dark:bg-gray-50 ">
											↵
										</span>
									</div>
								</Command.Empty>
								<div
									className="overflow-y-auto"
									ref={scrollContainerRef}
									key={'commands_div'}
								>
									<Command.List>
										{commandItemsAndGroups.map(({ groupLabel, items }) => {
											return (
												<div key={groupLabel}>
													<Command.Group
														heading={groupLabel}
														className="flex flex-col gap-2"
														key={groupLabel}
													>
														{items.map(
															({ icon, label, action, slug }) => {
																const Icon =
																	Icons[icon || 'arrowRight']

																const isSelected =
																	currentSelected === slug

																return (
																	<Command.Item
																		id={`command-item-${slug}`}
																		cmkd-item={slug}
																		key={slug}
																		value={slug}
																		cmdk-item=""
																		data-active={
																			isSelected
																				? 'true'
																				: 'false'
																		}
																		onSelect={value => {
																			runAction(action, value)
																		}}
																		className={clsx(
																			'flex items-center justify-between',
																			currentSelected === slug
																				? 'bg-accent text-accent-foreground'
																				: ''
																		)}
																	>
																		<div className="flex items-center gap-2">
																			<Icon className="size-4" />
																			{label}
																		</div>

																		{isSelected ? (
																			<div className="opacity-100 transition-all ease-in-out">
																				<span className="inline-flex h-5 w-fit items-center justify-center rounded-md px-2 py-0.5 text-center text-xs font-semibold leading-4 text-gray-700 dark:bg-gray-50 ">
																					↵
																				</span>
																			</div>
																		) : null}
																	</Command.Item>
																)
															}
														)}
													</Command.Group>
													<Command.Separator className="my-2 h-[1px] w-full bg-accent" />
												</div>
											)
										})}
									</Command.List>
								</div>
							</Command>
						</div>
					</motion.div>
				</DialogContent>
			</DialogPortal>
		</Dialog>
	)
}
