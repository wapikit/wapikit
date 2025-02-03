import * as React from 'react'
import { Slot } from '@radix-ui/react-slot'
import { cva, type VariantProps } from 'class-variance-authority'

import { clsx as cn } from 'clsx'

const buttonVariants = cva(
	'inline-flex min-w-max items-center justify-center flex-shrink-0 border font-normal focus:outline-none disabled:shadow-none disabled:cursor-not-allowed disabled:opacity-50 cursor-pointer',
	{
		variants: {
			variant: {
				default: 'bg-primary border text-primary-foreground shadow hover:bg-primary/90',
				destructive:
					'bg-destructive text-destructive-foreground shadow-sm hover:bg-destructive/90',
				outline:
					'border border-input bg-transparent shadow-sm hover:bg-accent hover:text-accent-foreground',
				text: 'border-none border-input bg-transparent shadow-sm hover:bg-accent hover:text-accent-foreground',
				secondary: 'bg-secondary text-secondary-foreground shadow-sm hover:bg-secondary/80',
				ghost: 'rounded-full border-none ',
				link: 'text-primary underline-offset-4 hover:underline'
			},
			size: {
				default: 'h-9 px-4 py-2 text-sm rounded-lg',
				xSmallForGraphics: 'rounded-[4px] px-1.5 py-[2px] text-[5px] gap-1',
				smallForGraphics: 'rounded-md px-1.5 py-1 text-[6px] gap-1',
				sm: 'rounded-lg px-2.5 py-1.5 text-xs gap-4 ',
				medium: 'rounded-lg px-3 py-1.5 text-sm gap-4 ',
				lg: 'rounded-lg px-4 py-2 text-sm gap-4 ',
				xLarge: 'rounded-lg px-6 py-1.5 text-base gap-4 ',
				badge: 'rounded-full px-3 py-1.5 text-xs gap-4 ',
				icon: 'rounded-full p-2'
			}
		},
		defaultVariants: {
			variant: 'default',
			size: 'default'
		}
	}
)

export interface ButtonProps
	extends React.ButtonHTMLAttributes<HTMLButtonElement>,
		VariantProps<typeof buttonVariants> {
	asChild?: boolean
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
	({ className, variant, size, asChild = false, ...props }, ref) => {
		const Comp = asChild ? Slot : 'button'
		return (
			<Comp
				className={cn(buttonVariants({ variant, size, className }))}
				ref={ref}
				{...props}
			/>
		)
	}
)
Button.displayName = 'Button'

export { Button, buttonVariants }
