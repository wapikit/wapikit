'use client'

import type { Attachment } from 'ai'
import { clsx as cx } from 'clsx'
import type React from 'react'
import { useRef, useEffect, useCallback, type Dispatch, type SetStateAction } from 'react'
import { toast } from 'sonner'
import { useLocalStorage, useWindowSize } from 'usehooks-ts'
import { PaperclipIcon, StopIcon } from './icons'
import { PreviewAttachment } from './preview-attachment'
import { Button } from '~/components/ui/button'
import { Textarea } from '../ui/textarea'
import { SuggestedActions } from './suggested-actions'
import { SendIcon } from 'lucide-react'
import { useAiChatStore } from '~/store/ai-chat-store'

const AiChatInput = ({
	isLoading,
	stop,
	attachments,
	setAttachments,
	handleSubmit,
	selectSuggestedAction,
	className
}: {
	chatId: string
	isLoading: boolean
	stop: () => void
	attachments: Array<Attachment>
	setAttachments: Dispatch<SetStateAction<Array<Attachment>>>
	selectSuggestedAction: (action: string) => void
	handleSubmit: (event?: { preventDefault?: () => void }) => void
	className?: string
}) => {
	const textareaRef = useRef<HTMLTextAreaElement>(null)
	const { width } = useWindowSize()

	const { writeProperty, inputValue } = useAiChatStore()

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

	const [localStorageInput, setLocalStorageInput] = useLocalStorage('input', '')

	useEffect(() => {
		if (textareaRef.current) {
			const domValue = textareaRef.current.value
			// Prefer DOM value over localStorage to handle hydration
			const finalValue = domValue || localStorageInput || ''
			writeProperty({
				inputValue: finalValue
			})
			adjustHeight()
		}
		// Only run once after hydration
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [])

	useEffect(() => {
		setLocalStorageInput(inputValue)
	}, [inputValue, setLocalStorageInput])

	const handleInput = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
		writeProperty({
			inputValue: event.target.value
		})
		adjustHeight()
	}

	const { currentChatMessages } = useAiChatStore()

	const submitForm = useCallback(() => {
		handleSubmit()
		setAttachments([])
		setLocalStorageInput('')

		if (width && width > 768) {
			textareaRef.current?.focus()
		}
	}, [handleSubmit, setAttachments, setLocalStorageInput, width])

	return (
		<div className="relative flex w-full flex-col gap-4">
			{currentChatMessages.length === 0 && attachments.length === 0 && (
				<SuggestedActions selectSuggestedAction={selectSuggestedAction} />
			)}

			{attachments.length > 0 && (
				<div className="flex flex-row items-end gap-2 overflow-x-scroll">
					{attachments.map(attachment => (
						<PreviewAttachment key={attachment.url} attachment={attachment} />
					))}
				</div>
			)}

			<Textarea
				ref={textareaRef}
				placeholder="Send a message..."
				value={inputValue}
				onChange={handleInput}
				className={cx(
					'relative max-h-[calc(75dvh)] min-h-[24px] resize-none overflow-hidden bg-muted pb-10 !text-base dark:border-zinc-700',
					className
				)}
				rows={2}
				autoFocus
				onKeyDown={event => {
					if (event.key === 'Enter' && !event.shiftKey) {
						event.preventDefault()

						if (isLoading) {
							toast.error('Please wait for the model to finish its response!')
						} else {
							submitForm()
						}
					}
				}}
			/>

			{/* <div className="absolute bottom-0 flex w-fit flex-row justify-start p-2">
				<AttachmentButton fileInputRef={fileInputRef} isLoading={isLoading} />
			</div> */}

			<div className="absolute bottom-0 right-0 flex w-fit flex-row justify-end p-2">
				{isLoading ? (
					<StopButton stop={stop} />
				) : (
					<SendButton input={inputValue} submitForm={submitForm} />
				)}
			</div>
		</div>
	)
}

const AttachmentButton = ({
	fileInputRef,
	isLoading
}: {
	fileInputRef: React.MutableRefObject<HTMLInputElement | null>
	isLoading: boolean
}) => {
	return (
		<Button
			className="h-fit rounded-md rounded-bl-lg p-[7px] hover:bg-zinc-200 dark:border-zinc-700 hover:dark:bg-zinc-900"
			onClick={event => {
				event.preventDefault()
				fileInputRef.current?.click()
			}}
			disabled={isLoading}
			variant="ghost"
		>
			<PaperclipIcon size={14} />
		</Button>
	)
}

const StopButton = ({ stop }: { stop: () => void }) => {
	return (
		<Button
			className="h-fit rounded-full border p-1.5 dark:border-zinc-600"
			onClick={event => {
				event.preventDefault()
				stop()

				// ! TODO: why do we even need this ???
				// setMessages(messages => sanitizeUIMessages(messages))
			}}
		>
			<StopIcon size={14} />
		</Button>
	)
}

const SendButton = ({ submitForm, input }: { submitForm: () => void; input: string }) => {
	return (
		<Button
			className="h-fit rounded-full border p-1.5 dark:border-zinc-600"
			onClick={event => {
				event.preventDefault()
				submitForm()
			}}
			disabled={input.length === 0}
		>
			<SendIcon size={14} />
		</Button>
	)
}

export { AiChatInput, SendButton, StopButton, AttachmentButton }
