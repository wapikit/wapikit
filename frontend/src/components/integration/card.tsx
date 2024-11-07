import Image from 'next/image'
import { Card } from '../ui/card'
import Link from 'next/link'
import { Badge } from '../ui/badge'

const IntegrationCard: React.FC<{
	slug: string
	name: string
	description: string
	icon: string
}> = ({ description, icon, name, slug }) => {
	return (
		<Card className="flex h-72 flex-1 flex-col items-center justify-between gap-4 p-5">
			<div className="flex flex-col items-center justify-start gap-3">
				<div className="flex w-full flex-row items-end justify-start gap-2">
					<span className="rounded-lg border p-1">
						<Image
							src={icon}
							height={500}
							width={500}
							className="aspect-square h-10 w-10"
							alt={`${name}-icon`}
						/>
					</span>
					<h2 className="pb-0.5 text-xl font-bold">{name}</h2>
				</div>

				<p className="text-sm text-foreground">{description}</p>
			</div>

			<Link href={`/integrations/${slug}`} className="mr-auto">
				{/* <Button>Install</Button> */}
				<Badge>Coming Soon</Badge>
			</Link>
		</Card>
	)
}

export default IntegrationCard
