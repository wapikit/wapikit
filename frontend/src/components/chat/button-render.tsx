import Link from 'next/link'
import React from 'react'
import { type TemplateMessageComponentButton } from 'root/.generated'
import { Icons } from '../icons'

const MessageButtonRenderer: React.FC<{
	messageButton: TemplateMessageComponentButton
}> = ({ messageButton }) => {
	const UrlIcon = Icons.externalLink
	const PhoneIcon = Icons.phone
	const ClipBoardIcon = Icons.clipboard
	const ReplyIcon = Icons.reply

	if (messageButton.type === 'URL') {
		return (
			<Link href={messageButton.url || '/'}>
				<div className="flex items-center justify-center gap-2 py-2 text-center text-blue-500">
					<UrlIcon className="size-5" />
					{messageButton.text}
				</div>
			</Link>
		)
	} else if (messageButton.type === 'PHONE_NUMBER') {
		return (
			<div className="flex cursor-pointer items-center justify-center gap-2 py-2 text-center text-blue-500">
				<PhoneIcon className="size-5" />
				{messageButton.text}
			</div>
		)
	} else if (messageButton.type === 'COPY_CODE') {
		return (
			<div className="flex cursor-pointer items-center  justify-center gap-2 py-2 text-center text-blue-500">
				<ClipBoardIcon className="size-5" />
				{messageButton.text}
			</div>
		)
	} else {
		return (
			<div className="flex cursor-pointer items-center  justify-center gap-2 py-2 text-center text-blue-500">
				<ReplyIcon className="size-4" />
				{messageButton.text}
			</div>
		)
	}
}

export default MessageButtonRenderer
