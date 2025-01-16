'use client'

import { Button } from '~/components/ui/button'
import { type Dispatch, type SetStateAction, useEffect, useRef, useState } from 'react'
import { Textarea } from '../ui/textarea'
import { toast } from 'sonner'
import { useAiChatStore } from '~/store/ai-chat-store'

export type MessageEditorProps = {
	messageId: string
	setMode: Dispatch<SetStateAction<'view' | 'edit'>>
}

export function MessageEditor({ messageId, setMode }: MessageEditorProps) {
	const [isSubmitting, setIsSubmitting] = useState<boolean>(false)
	const { currentChatMessages, editMessage } = useAiChatStore()

	const message = currentChatMessages.find(m => m.uniqueId === messageId)

	const [draftContent, setDraftContent] = useState<string>(message?.content || '')

	const textareaRef = useRef<HTMLTextAreaElement>(null)

	useEffect(() => {
		if (textareaRef.current) {
			adjustHeight()
		}
	}, [])

	const adjustHeight = () => {
		if (textareaRef.current) {
			textareaRef.current.style.height = 'auto'
			textareaRef.current.style.height = `${textareaRef.current.scrollHeight + 2}px`
		}
	}

	const handleInput = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
		setDraftContent(event.target.value)
		adjustHeight()
	}

	return (
		<div className="flex w-full flex-col gap-2">
			<Textarea
				ref={textareaRef}
				className="w-full resize-none overflow-hidden rounded-md bg-transparent !text-base outline-none"
				value={draftContent}
				onChange={handleInput}
			/>

			<div className="flex flex-row justify-end gap-2">
				<Button
					variant="outline"
					className="h-fit px-3 py-2"
					onClick={() => {
						setMode('view')
					}}
				>
					Cancel
				</Button>
				<Button
					variant="default"
					className="h-fit px-3 py-2"
					disabled={isSubmitting}
					onClick={async () => {
						setIsSubmitting(true)
						const messageId = message?.uniqueId

						if (!messageId) {
							toast.error('Something went wrong, please try again!')
							setIsSubmitting(false)
							return
						}

						// ! TODO: AI: Implement deleteTrailingMessages
						// await deleteTrailingMessages({
						// 	id: messageId
						// })

						editMessage(message.uniqueId, draftContent)
						setMode('view')

						// ! TODO: AI: Implement deleteTrailingMessages
					}}
				>
					{isSubmitting ? 'Sending...' : 'Send'}
				</Button>
			</div>
		</div>
	)
}
