'use client'

import { Command } from 'cmdk'
import { useEffect } from 'react'
import { useLayoutStore } from '~/store/layout.store'
import { motion } from 'framer-motion'
import { Popover, PopoverContent, PopoverTrigger } from '@radix-ui/react-popover'

export default function CommandMenuProvider() {
	const { isCommandMenuOpen, writeProperty } = useLayoutStore()

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
		}

		document.addEventListener('keydown', down)
		return () => document.removeEventListener('keydown', down)
	}, [writeProperty, isCommandMenuOpen])

	return (
		<Popover
			onOpenChange={isOpen => {
				writeProperty({
					isCommandMenuOpen: isOpen
				})
			}}
			open={isCommandMenuOpen}
			key={'command_menu'}
		>
			<PopoverTrigger>Toggle popover</PopoverTrigger>
			<PopoverContent>
				<motion.div
					initial={{ opacity: 0, scale: 0.98 }}
					animate={{ opacity: 1, scale: 1 }}
					exit={{ opacity: 0, scale: 0.98 }}
					transition={{ duration: 0.2 }}
					style={{
						height: 475
					}}
				>
					<div className="linear relative z-[99999]">
						<Command className="">
							{/* <Command cmdk-linear-badge="">Issue - FUN-343</Command> */}
							<Command.Input autoFocus placeholder="Type a command or search..." />

							<Command.Group heading="Campaigns">
								<Command.List>
									<Command.Empty>No results found.</Command.Empty>
									{campaignItems.map(({ icon, label }) => {
										return (
											<Command.Item key={label} value={label} cmdk-item="">
												{icon}
												{label}
											</Command.Item>
										)
									})}
								</Command.List>
							</Command.Group>

							<Command.Separator />

							<Command.Group heading="Contacts">
								<Command.List>
									<Command.Empty>No results found.</Command.Empty>
									{contactItems.map(({ icon, label }) => {
										return (
											<Command.Item key={label} value={label}>
												{icon}
												{label}
											</Command.Item>
										)
									})}
								</Command.List>
							</Command.Group>
						</Command>
					</div>
				</motion.div>
			</PopoverContent>
		</Popover>
	)
}

const campaignItems = [
	{
		icon: <AssignToIcon />,
		label: 'Create Campaign...'
	}
]

const contactItems = [
	{
		icon: <AssignToIcon />,
		label: 'Create List...'
	},
	{
		icon: <AssignToMeIcon />,
		label: 'Bulk Import Contacts...'
	},
	{
		icon: <AssignToMeIcon />,
		label: 'Create Contact...'
	}
]

function AssignToIcon() {
	return (
		<svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
			<path d="M7 7a2.5 2.5 0 10.001-4.999A2.5 2.5 0 007 7zm0 1c-1.335 0-4 .893-4 2.667v.666c0 .367.225.667.5.667h2.049c.904-.909 2.417-1.911 4.727-2.009v-.72a.27.27 0 01.007-.063C9.397 8.404 7.898 8 7 8zm4.427 2.028a.266.266 0 01.286.032l2.163 1.723a.271.271 0 01.013.412l-2.163 1.97a.27.27 0 01-.452-.2v-.956c-3.328.133-5.282 1.508-5.287 1.535a.27.27 0 01-.266.227h-.022a.27.27 0 01-.249-.271c0-.046 1.549-3.328 5.824-3.509v-.72a.27.27 0 01.153-.243z" />
		</svg>
	)
}

function AssignToMeIcon() {
	return (
		<svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
			<path d="M7.00003 7C8.38128 7 9.50003 5.88125 9.50003 4.5C9.50003 3.11875 8.38128 2 7.00003 2C5.61878 2 4.50003 3.11875 4.50003 4.5C4.50003 5.88125 5.61878 7 7.00003 7Z" />
			<path
				fillRule="evenodd"
				clipRule="evenodd"
				d="M7.00005 8C5.66505 8 3.00006 8.89333 3.00006 10.6667V11.3333C3.00006 11.7 3.22506 12 3.50006 12H3.98973C4.01095 11.9415 4.04535 11.8873 4.09266 11.8425L7.21783 8.88444C7.28966 8.81658 7.38297 8.77917 7.4796 8.77949C7.69459 8.78018 7.86826 8.96356 7.86753 9.1891L7.86214 10.629C9.00553 10.5858 10.0366 10.4354 10.9441 10.231C10.5539 8.74706 8.22087 8 7.00005 8Z"
			/>
			<path d="M6.72511 14.718C6.80609 14.7834 6.91767 14.7955 7.01074 14.749C7.10407 14.7036 7.16321 14.6087 7.16295 14.5047L7.1605 13.7849C11.4352 13.5894 12.9723 10.3023 12.9722 10.2563C12.9722 10.1147 12.8634 9.9971 12.7225 9.98626L12.7009 9.98634C12.5685 9.98689 12.4561 10.0833 12.4351 10.2142C12.4303 10.2413 10.4816 11.623 7.15364 11.7666L7.1504 10.8116C7.14981 10.662 7.02829 10.5412 6.87896 10.5418C6.81184 10.5421 6.74721 10.5674 6.69765 10.6127L4.54129 12.5896C4.43117 12.6906 4.42367 12.862 4.52453 12.9723C4.53428 12.9829 4.54488 12.9928 4.55621 13.0018L6.72511 14.718Z" />
		</svg>
	)
}
