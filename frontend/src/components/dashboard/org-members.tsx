import { Avatar, AvatarFallback, AvatarImage } from '~/components/ui/avatar'

export function OrganizationMembers() {
	const members = [
		{
			name: 'Olivia Martin',
			email: 'olivia@gmail.com',
			lastActive: 'Active'
		},
		{
			name: 'Jackson Lee',
			email: 'jackson@gmail.com',
			lastActive: '3 hours ago'
		},
		{
			name: 'Isabella Nguyen',
			email: 'isabella@gmail.com',
			lastActive: '1 day ago'
		},
		{
			name: 'Aiden Smith',
			email: 'aiden@gmail.com',
			lastActive: '1 week ago'
		}
	]

	return (
		<div className="space-y-8 px-6">
			{members.map((member, index) => {
				return (
					<div className="flex items-center" key={index}>
						<Avatar className="h-9 w-9">
							<AvatarImage src="/avatars/01.png" alt="Avatar" />
							<AvatarFallback>OM</AvatarFallback>
						</Avatar>
						<div className="ml-4 space-y-1">
							<p className="text-sm font-medium leading-none">{member.name}</p>
							<p className="text-sm text-muted-foreground">{member.email}</p>
						</div>
						<div className="ml-auto text-sm font-normal">{member.lastActive}</div>
					</div>
				)
			})}
		</div>
	)
}
