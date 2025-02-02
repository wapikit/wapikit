'use client'

import { clsx as cx } from 'clsx'
import { AnimatePresence, motion } from 'framer-motion'
import { useEffect, useRef, useState } from 'react'
import { PencilEditIcon, SparklesIcon } from './icons'
import { Markdown } from './markdown'
import { MessageActions } from './message-actions'
import { clsx } from 'clsx'
import { Button } from '~/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger, TooltipProvider } from '../ui/tooltip'
import { MessageEditor } from './message-editor'
import {
	AiChatMessageRoleEnum,
	type AiChatMessageSchema,
	type AiChatMessageVoteSchema
} from 'root/.generated'

const PreviewMessage = ({
	chatId,
	message,
	vote,
	isLoading,
	isReadonly,
	isCurrentMessageInProgress
}: {
	chatId: string
	message: AiChatMessageSchema
	vote: AiChatMessageVoteSchema | undefined
	isLoading: boolean
	isReadonly: boolean
	isCurrentMessageInProgress: boolean
}) => {
	const [mode, setMode] = useState<'view' | 'edit'>('view')
	const [typedContent, setTypedContent] = useState('')
	const typingIndex = useRef(0)
	const typingInterval = useRef<NodeJS.Timeout | null>(null)

	useEffect(() => {
		// Only apply typing animation for the current message in progress
		if (isCurrentMessageInProgress) {
			typingIndex.current = 0
			setTypedContent('')

			typingInterval.current = setInterval(() => {
				if (typingIndex.current < message.content.length) {
					setTypedContent(prev => prev + message.content[typingIndex.current])
					typingIndex.current++
				} else {
					if (typingInterval.current) clearInterval(typingInterval.current)
				}
			}, 100) // Adjust typing speed here
		}

		return () => {
			if (typingInterval.current) clearInterval(typingInterval.current)
		}
	}, [isCurrentMessageInProgress, message.content])

	return (
		<AnimatePresence>
			<motion.div
				className="group/message mx-auto w-full px-4"
				initial={{ y: 5, opacity: 0 }}
				animate={{ y: 0, opacity: 1 }}
				data-role={message.role}
			>
				<div
					className={clsx(
						'flex w-full gap-4 group-data-[role=user]/message:ml-auto group-data-[role=user]/message:max-w-2xl',
						{
							'w-full': mode === 'edit',
							'group-data-[role=user]/message:w-fit': mode !== 'edit'
						}
					)}
				>
					{message.role === AiChatMessageRoleEnum.Assistant && (
						<div className="flex size-6 shrink-0 items-center justify-center rounded-full bg-background ring-1 ring-border">
							<div className="translate-y-px">
								<SparklesIcon size={12} />
							</div>
						</div>
					)}

					<div className="flex w-full flex-col gap-2">
						{/* {message.experimental_attachments && (
							<div className="flex flex-row justify-end gap-2">
								{message.experimental_attachments.map(attachment => (
									<PreviewAttachment
										key={attachment.url}
										attachment={attachment}
									/>
								))}
							</div>
						)} */}

						{message.content && mode === 'view' && (
							<div
								className={clsx(
									'flex flex-row items-center gap-2',
									message.role === AiChatMessageRoleEnum.User && 'justify-end'
								)}
							>
								{message.role === AiChatMessageRoleEnum.User && !isReadonly && (
									<TooltipProvider>
										<Tooltip>
											<TooltipTrigger asChild>
												<Button
													variant="ghost"
													className="h-fit rounded-full px-2 text-muted-foreground opacity-0 group-hover/message:opacity-100"
													onClick={() => {
														setMode('edit')
													}}
												>
													<PencilEditIcon />
												</Button>
											</TooltipTrigger>
											<TooltipContent>Edit message</TooltipContent>
										</Tooltip>
									</TooltipProvider>
								)}

								<div
									className={clsx(
										'flex max-w-lg flex-col items-center justify-start gap-4 px-3 py-1 text-sm',
										message.role === AiChatMessageRoleEnum.User
											? 'rounded-md bg-primary text-primary-foreground'
											: ''
									)}
								>
									{/* Typing animation for last message */}
									{isCurrentMessageInProgress ? (
										<div>
											<Markdown>{typedContent}</Markdown>
											<span className="blinking-cursor">|</span>
										</div>
									) : (
										<Markdown>{message.content}</Markdown>
									)}
								</div>
							</div>
						)}

						{message.content && mode === 'edit' && (
							<div className="flex flex-row items-start gap-2">
								<div className="size-8" />
								<MessageEditor
									key={message.uniqueId}
									setMode={setMode}
									messageId={message.uniqueId}
								/>
							</div>
						)}

						{!isReadonly && (
							<MessageActions
								key={`action-${message.uniqueId}`}
								chatId={chatId}
								message={message}
								vote={vote}
								isLoading={isLoading}
							/>
						)}
					</div>
				</div>
			</motion.div>
		</AnimatePresence>
	)
}

const ThinkingMessage = () => {
	const role = 'assistant'

	return (
		<motion.div
			className="group/message mx-auto w-full px-4 "
			initial={{ y: 5, opacity: 0 }}
			animate={{ y: 0, opacity: 1, transition: { delay: 1 } }}
			data-role={role}
		>
			<div
				className={cx(
					'flex w-full gap-4 rounded-xl group-data-[role=user]/message:ml-auto group-data-[role=user]/message:w-fit group-data-[role=user]/message:max-w-2xl group-data-[role=user]/message:px-3 group-data-[role=user]/message:py-2',
					{
						'group-data-[role=user]/message:bg-muted': true
					}
				)}
			>
				<div className="flex size-8 shrink-0 items-center justify-center rounded-full ring-1 ring-border">
					<SparklesIcon size={14} />
				</div>

				<div className="flex w-full flex-col gap-2">
					<div className="flex flex-col gap-4 text-muted-foreground">Thinking...</div>
				</div>
			</div>
		</motion.div>
	)
}

export { PreviewMessage, ThinkingMessage }
