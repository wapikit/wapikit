'use client'

import { useRouter } from 'next/navigation'
import { Avatar, AvatarFallback, AvatarImage } from '~/components/ui/avatar'
import { Button } from '~/components/ui/button'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuGroup,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuTrigger
} from '~/components/ui/dropdown-menu'
import { useAuthState } from '~/hooks/use-auth-state'

export function UserNav() {
	const router = useRouter()
	const { authState } = useAuthState()
	if (authState.isAuthenticated) {
		return (
			<DropdownMenu>
				<DropdownMenuTrigger asChild>
					<Button variant="ghost" size={'icon'} className="relative">
						<Avatar className="h-8 w-8">
							<AvatarImage
								src={'/assets/empty-pfp.png'}
								alt={authState.data.user.name}
							/>
							<AvatarFallback>{authState.data.user.name}</AvatarFallback>
						</Avatar>
					</Button>
				</DropdownMenuTrigger>
				<DropdownMenuContent className="w-56" align="end" forceMount>
					<DropdownMenuLabel className="font-normal">
						<div className="flex flex-col space-y-1">
							<p className="text-sm font-medium leading-none">
								{authState.data.user.name}
							</p>
							<p className="text-xs leading-none text-muted-foreground">
								{authState.data.user.email}
							</p>
						</div>
					</DropdownMenuLabel>
					<DropdownMenuSeparator />
					<DropdownMenuGroup>
						<DropdownMenuItem
							onClick={() => {
								router.push('/settings')
							}}
						>
							Settings
						</DropdownMenuItem>
						<DropdownMenuItem
							onClick={() => {
								router.push('/settings?tab=api-key')
							}}
						>
							API Access
						</DropdownMenuItem>
						<DropdownMenuItem
							onClick={() => {
								window.open('https://docs.wapikit.com', '_blank')
							}}
						>
							Documentation
						</DropdownMenuItem>
					</DropdownMenuGroup>
					<DropdownMenuSeparator />
					<DropdownMenuItem
						onClick={() => {
							router.push('/logout')
						}}
					>
						Log out
					</DropdownMenuItem>
				</DropdownMenuContent>
			</DropdownMenu>
		)
	} else {
		return null
	}
}
