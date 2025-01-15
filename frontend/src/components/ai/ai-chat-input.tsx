'use client'

import type { Attachment } from 'ai'
import { clsx as cx } from 'clsx'
import type React from 'react'
import {
	useRef,
	useEffect,
	useState,
	useCallback,
	type Dispatch,
	type SetStateAction,
	type ChangeEvent
} from 'react'
import { toast } from 'sonner'
import { useLocalStorage, useWindowSize } from 'usehooks-ts'
import { sanitizeUIMessages } from '~/utils/ai-utils'
import { PaperclipIcon, StopIcon } from './icons'
import { PreviewAttachment } from './preview-attachment'
import { Button } from '~/components/ui/button'
import { Textarea } from '../ui/textarea'
import { SuggestedActions } from './suggested-actions'
import { SendIcon } from 'lucide-react'
import { useAiChatStore } from '~/store/ai-chat-store'

const AiChatInput = ({
	chatId,
	input,
	setInput,
	isLoading,
	stop,
	attachments,
	setAttachments,
	handleSubmit,
	selectSuggestedAction,
	className
}: {
	chatId: string
	input: string
	setInput: (value: string) => void
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
			setInput(finalValue)
			adjustHeight()
		}
		// Only run once after hydration
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [])

	useEffect(() => {
		setLocalStorageInput(input)
	}, [input, setLocalStorageInput])

	const handleInput = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
		setInput(event.target.value)
		adjustHeight()
	}

	const fileInputRef = useRef<HTMLInputElement>(null)
	const [uploadQueue, setUploadQueue] = useState<Array<string>>([])

	const { currentChatMessages } = useAiChatStore()

	const submitForm = useCallback(() => {
		handleSubmit()
		setAttachments([])
		setLocalStorageInput('')

		if (width && width > 768) {
			textareaRef.current?.focus()
		}
	}, [attachments, handleSubmit, setAttachments, setLocalStorageInput, width, chatId])

	const uploadFile = async (file: File) => {
		const formData = new FormData()
		formData.append('file', file)

		try {
			const response = await fetch('/api/files/upload', {
				method: 'POST',
				body: formData
			})

			if (response.ok) {
				const data = await response.json()
				const { url, pathname, contentType } = data

				return {
					url,
					name: pathname,
					contentType: contentType
				}
			}
			const { error } = await response.json()
			toast.error(error)
		} catch (error) {
			toast.error('Failed to upload file, please try again!')
		}
	}

	const handleFileChange = useCallback(
		async (event: ChangeEvent<HTMLInputElement>) => {
			const files = Array.from(event.target.files || [])

			setUploadQueue(files.map(file => file.name))

			try {
				const uploadPromises = files.map(file => uploadFile(file))
				const uploadedAttachments = await Promise.all(uploadPromises)
				const successfullyUploadedAttachments = uploadedAttachments.filter(
					attachment => attachment !== undefined
				)

				setAttachments(currentAttachments => [
					...currentAttachments,
					...successfullyUploadedAttachments
				])
			} catch (error) {
				console.error('Error uploading files!', error)
			} finally {
				setUploadQueue([])
			}
		},
		[setAttachments]
	)

	return (
		<div className="relative flex w-full flex-col gap-4">
			{currentChatMessages.length === 0 &&
				attachments.length === 0 &&
				uploadQueue.length === 0 && (
					<SuggestedActions selectSuggestedAction={selectSuggestedAction} />
				)}

			<input
				type="file"
				className="pointer-events-none fixed -left-4 -top-4 size-0.5 opacity-0"
				ref={fileInputRef}
				multiple
				onChange={handleFileChange}
				tabIndex={-1}
			/>

			{(attachments.length > 0 || uploadQueue.length > 0) && (
				<div className="flex flex-row items-end gap-2 overflow-x-scroll">
					{attachments.map(attachment => (
						<PreviewAttachment key={attachment.url} attachment={attachment} />
					))}

					{uploadQueue.map(filename => (
						<PreviewAttachment
							key={filename}
							attachment={{
								url: '',
								name: filename,
								contentType: ''
							}}
							isUploading={true}
						/>
					))}
				</div>
			)}

			<Textarea
				ref={textareaRef}
				placeholder="Send a message..."
				value={input}
				onChange={handleInput}
				className={cx(
					'max-h-[calc(75dvh)] min-h-[24px] resize-none overflow-hidden bg-muted pb-10 !text-base dark:border-zinc-700',
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

			<div className="absolute bottom-0 flex w-fit flex-row justify-start p-2">
				<AttachmentButton fileInputRef={fileInputRef} isLoading={isLoading} />
			</div>

			<div className="absolute bottom-0 right-0 flex w-fit flex-row justify-end p-2">
				{isLoading ? (
					<StopButton stop={stop} />
				) : (
					<SendButton input={input} submitForm={submitForm} uploadQueue={uploadQueue} />
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

const SendButton = ({
	submitForm,
	input,
	uploadQueue
}: {
	submitForm: () => void
	input: string
	uploadQueue: Array<string>
}) => {
	return (
		<Button
			className="h-fit rounded-full border p-1.5 dark:border-zinc-600"
			onClick={event => {
				event.preventDefault()
				submitForm()
			}}
			disabled={input.length === 0 || uploadQueue.length > 0}
		>
			<SendIcon size={14} />
		</Button>
	)
}

export { AiChatInput, SendButton, StopButton, AttachmentButton }
