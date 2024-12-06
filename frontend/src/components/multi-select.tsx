'use client'

import * as React from 'react'
import { cva, type VariantProps } from 'class-variance-authority'
import { XCircle, XIcon } from 'lucide-react'
import { clsx } from 'clsx'
import { useState } from 'react'
import { Popover, PopoverContent, PopoverTrigger } from './ui/popover'
import { Badge } from './ui/badge'
import { Separator } from './ui/separator'
import {
	Command,
	CommandEmpty,
	CommandGroup,
	CommandInput,
	CommandItem,
	CommandList,
	CommandSeparator
} from './ui/command'
import { CaretSortIcon } from '@radix-ui/react-icons'

const multiSelectVariants = cva(
	'm-1 transition ease-in-out delay-150 hover:-translate-y-1 hover:scale-110 duration-300',
	{
		variants: {
			variant: {
				default: 'border-foreground/10 text-gray-500 bg-card hover:bg-card/80',
				secondary:
					'border-foreground/10 bg-secondary text-secondary-foreground hover:bg-secondary/80',
				destructive:
					'border-transparent bg-destructive text-destructive-foreground hover:bg-destructive/80',
				inverted: 'inverted'
			}
		},
		defaultVariants: {
			variant: 'default'
		}
	}
)

interface MultiSelectProps
	extends React.HTMLAttributes<HTMLDivElement>,
		VariantProps<typeof multiSelectVariants> {
	options: {
		value: string
		label?: string
	}[]
	filter?: boolean
	onValueChange: (value: string[]) => void
	defaultValue: string[]
	placeholder?: string
	maxCount?: number
	modalPopover?: boolean
	// asChild?: boolean
	className?: string
}
const CheckFilled = () => {
	return (
		<svg
			xmlns="http://www.w3.org/2000/svg"
			width="18"
			height="18"
			viewBox="0 0 18 18"
			fill="#202020"
			className={clsx('h-8 w-8 ')}
		>
			<path
				d="M5.25 3.75C4.85218 3.75 4.47064 3.90804 4.18934 4.18934C3.90804 4.47064 3.75 4.85218 3.75 5.25V12.75C3.75 13.1478 3.90804 13.5294 4.18934 13.8107C4.47064 14.092 4.85218 14.25 5.25 14.25H12.75C13.1478 14.25 13.5294 14.092 13.8107 13.8107C14.092 13.5294 14.25 13.1478 14.25 12.75V5.25C14.25 4.85218 14.092 4.47064 13.8107 4.18934C13.5294 3.90804 13.1478 3.75 12.75 3.75H5.25ZM8.25 11.5605L6.21975 9.53025L7.28025 8.46975L8.25 9.4395L11.0948 6.59475L12.1553 7.65525L8.25 11.5605Z"
				fill="#ffffff"
			/>
		</svg>
	)
}

export const MultiSelect = React.forwardRef<HTMLButtonElement, MultiSelectProps>(
	(
		{
			filter = false,
			options,
			onValueChange,
			variant,
			defaultValue = [],
			placeholder = 'Select options',
			maxCount = 3,
			modalPopover = false,
			// asChild = false,
			className,
			...props
		},
		ref
	) => {
		const [selectedValues, setSelectedValues] = useState<string[]>(defaultValue)

		const [isPopoverOpen, setIsPopoverOpen] = useState(false)

		const handleInputKeyDown = (event: React.KeyboardEvent<HTMLInputElement>) => {
			if (event.key === 'Enter') {
				setIsPopoverOpen(true)
			} else if (event.key === 'Backspace' && !event.currentTarget.value) {
				const newSelectedValues = [...selectedValues]
				newSelectedValues.pop()
				setSelectedValues(newSelectedValues)
				onValueChange(newSelectedValues)
			}
		}

		const toggleOption = (value: string) => {
			const newSelectedValues = selectedValues.includes(value)
				? selectedValues.filter(v => v !== value)
				: [...selectedValues, value]
			setSelectedValues(newSelectedValues)
			onValueChange(newSelectedValues)
		}

		const handleClear = () => {
			setSelectedValues([])
			onValueChange([])
		}

		const handleTogglePopover = () => {
			setIsPopoverOpen(prev => !prev)
		}

		const clearExtraOptions = () => {
			const newSelectedValues = selectedValues.slice(0, maxCount)
			setSelectedValues(newSelectedValues)
			onValueChange(newSelectedValues)
		}

		const toggleAll = () => {
			if (selectedValues.length === options.length) {
				handleClear()
			} else {
				const allValues = options.map(option => option.value)
				setSelectedValues(allValues)
				onValueChange(allValues)
			}
		}

		return (
			<Popover open={isPopoverOpen} onOpenChange={setIsPopoverOpen} modal={modalPopover}>
				<PopoverTrigger>
					<div
						ref={ref as any}
						{...props}
						onClick={handleTogglePopover}
						className={clsx(
							'flex w-full cursor-pointer items-center justify-between rounded-md border border-input bg-inherit px-3 py-2 text-sm font-medium shadow-sm hover:bg-inherit',
							className
						)}
					>
						{!filter ? (
							<>
								{selectedValues.length > 0 ? (
									<div className="flex w-full items-center justify-between">
										<div className="flex flex-wrap items-center">
											{selectedValues.slice(0, maxCount).map(value => {
												const option = options.find(o => o.value === value)
												return (
													<Badge
														key={value}
														className={clsx(
															' py-0.5',
															multiSelectVariants({ variant })
														)}
													>
														{option?.label || option?.value}
														<XCircle
															className="ml-2 h-4 w-4 cursor-pointer"
															onClick={event => {
																event.stopPropagation()
																toggleOption(value)
															}}
														/>
													</Badge>
												)
											})}
											{selectedValues.length > maxCount && (
												<Badge
													className={clsx(
														'border-foreground/1 bg-transparent text-foreground hover:bg-transparent',
														multiSelectVariants({ variant })
													)}
												>
													{`+ ${selectedValues.length - maxCount} more`}
													<XCircle
														className="ml-2 h-4 w-4 cursor-pointer"
														onClick={event => {
															event.stopPropagation()
															clearExtraOptions()
														}}
													/>
												</Badge>
											)}
										</div>
										<div className="flex items-center justify-between">
											<XIcon
												className="mx-2 h-4 w-4 cursor-pointer opacity-50"
												onClick={event => {
													event.stopPropagation()
													handleClear()
												}}
											/>
											<Separator
												orientation="vertical"
												className="mr-2 flex h-full min-h-6"
											/>
											<CaretSortIcon className="h-4 w-4 opacity-50" />
										</div>
									</div>
								) : (
									<div className="mx-auto flex w-full items-center justify-between">
										<span className="mx-2 text-sm text-muted-foreground">
											{placeholder}
										</span>
										<CaretSortIcon className="h-4 w-4 opacity-50" />
									</div>
								)}
							</>
						) : (
							<div className="mx-auto flex w-full items-center justify-between">
								<span className="mx-2 text-sm text-gray-300">{placeholder}</span>
								<CaretSortIcon className="h-4 w-4 opacity-50" />
							</div>
						)}
					</div>
				</PopoverTrigger>
				<PopoverContent
					className="w-auto p-0"
					align="start"
					onEscapeKeyDown={() => setIsPopoverOpen(false)}
				>
					<Command className="min-w-44">
						{options.length > 0 ? (
							<CommandInput
								placeholder="Search..."
								onKeyDown={handleInputKeyDown}
								className="my-2 rounded-md !p-1 !pl-2 focus:border-gray-200"
							/>
						) : (
							<CommandItem>No options to select</CommandItem>
						)}
						<CommandList>
							<CommandEmpty>No results found.</CommandEmpty>
							<CommandGroup>
								{options.length > 0 && (
									<CommandItem
										key="all"
										onSelect={toggleAll}
										className="cursor-pointer"
									>
										<div
											className={clsx(
												'mr-2 flex h-4 w-4 items-center justify-center rounded-sm border border-primary',
												selectedValues.length === options.length
													? 'bg-primary text-primary-foreground'
													: 'opacity-50 [&_svg]:invisible'
											)}
										>
											<CheckFilled />
										</div>
										<span>(Select All)</span>
									</CommandItem>
								)}
								{options.map(option => {
									const isSelected = selectedValues.includes(option.value)
									return (
										<CommandItem
											key={option.value}
											onSelect={() => toggleOption(option.value)}
											className="cursor-pointer"
										>
											<div
												className={clsx(
													'mr-2 flex h-4 w-4 items-center justify-center rounded-sm border border-primary',
													isSelected
														? 'bg-primary text-primary-foreground'
														: 'opacity-50 [&_svg]:invisible'
												)}
											>
												<CheckFilled />
											</div>
											<span>{option?.label || option.value}</span>
										</CommandItem>
									)
								})}
							</CommandGroup>
							<CommandSeparator />
							<CommandGroup>
								<div className="flex items-center justify-between">
									{selectedValues.length > 0 && (
										<>
											<CommandItem
												onSelect={handleClear}
												className="flex-1 cursor-pointer justify-center"
											>
												Clear
											</CommandItem>
											<Separator
												orientation="vertical"
												className="flex h-full min-h-6"
											/>
										</>
									)}
									<CommandItem
										onSelect={() => setIsPopoverOpen(false)}
										className="max-w-full flex-1 cursor-pointer justify-center"
									>
										Close
									</CommandItem>
								</div>
							</CommandGroup>
						</CommandList>
					</Command>
				</PopoverContent>
			</Popover>
		)
	}
)

MultiSelect.displayName = 'MultiSelect'
