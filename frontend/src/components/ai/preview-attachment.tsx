import type { Attachment } from 'ai'

import { LoaderIcon } from './icons'

export const PreviewAttachment = ({
	attachment,
	isUploading = false
}: {
	attachment: Attachment
	isUploading?: boolean
}) => {
	const { name, url, contentType } = attachment

	return (
		<div className="flex flex-col gap-2">
			<div className="relative flex aspect-video h-16 w-20 flex-col items-center justify-center rounded-md bg-muted">
				{contentType ? (
					contentType.startsWith('image') ? (
						// NOTE: it is recommended to use next/image for images
						// eslint-disable-next-line @next/next/no-img-element
						<img
							key={url}
							src={url}
							alt={name ?? 'An image attachment'}
							className="size-full rounded-md object-cover"
						/>
					) : (
						<div className="" />
					)
				) : (
					<div className="" />
				)}

				{isUploading && (
					<div className="absolute animate-spin text-zinc-500">
						<LoaderIcon />
					</div>
				)}
			</div>
			<div className="max-w-16 truncate text-xs text-zinc-500">{name}</div>
		</div>
	)
}
