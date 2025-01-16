import { useCopyToClipboard } from 'usehooks-ts'
import { CopyIcon, ThumbDownIcon, ThumbUpIcon } from './icons'
import { Button } from '~/components/ui/button'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '../ui/tooltip'
import {
	AiChatMessageRoleEnum,
	type AiChatMessageSchema,
	useVoteOnAiChatMessage,
	type AiChatMessageVoteSchema
} from 'root/.generated'
import { successNotification } from '~/reusable-functions'

const MessageActions = ({
	chatId,
	message,
	vote,
	isLoading
}: {
	chatId: string
	message: AiChatMessageSchema
	vote: AiChatMessageVoteSchema | undefined
	isLoading: boolean
}) => {
	const copyToClipboard = useCopyToClipboard()[1]
	const voteMutation = useVoteOnAiChatMessage()

	if (isLoading) return null
	if (message.role === AiChatMessageRoleEnum.User) return null
	// if (message.toolInvocations && message.toolInvocations.length > 0) return null

	return (
		<TooltipProvider delayDuration={0}>
			<div className="flex flex-row gap-2">
				<Tooltip>
					<TooltipTrigger asChild>
						<Button
							className="h-fit text-muted-foreground"
							variant="outline"
							onClick={() => {
								copyToClipboard(message.content).catch(error =>
									console.error(error)
								)
								successNotification({
									message: 'Copied to clipboard!'
								})
							}}
						>
							<CopyIcon />
						</Button>
					</TooltipTrigger>
					<TooltipContent>Copy</TooltipContent>
				</Tooltip>

				<Tooltip>
					<TooltipTrigger asChild>
						<Button
							className="!pointer-events-auto h-fit px-2 py-1 text-muted-foreground"
							disabled={vote?.vote === 'Upvote'}
							variant="outline"
							onClick={() => {
								voteMutation
									.mutateAsync(
										{
											data: {
												messageId: message.uniqueId,
												type: 'Upvote'
											},
											id: chatId
										},
										{
											onSuccess: data => {
												if (data.vote) {
													successNotification({
														message: 'Voted!'
													})
												}
											}
										}
									)
									.catch(error => {
										console.error(error)
									})
							}}
						>
							<ThumbUpIcon />
						</Button>
					</TooltipTrigger>
					<TooltipContent>Upvote Response</TooltipContent>
				</Tooltip>

				<Tooltip>
					<TooltipTrigger asChild>
						<Button
							className="!pointer-events-auto h-fit px-2 py-1 text-muted-foreground"
							variant="outline"
							disabled={vote && vote.vote === 'Downvote'}
							onClick={() => {
								voteMutation
									.mutateAsync(
										{
											data: {
												messageId: message.uniqueId,
												type: 'Downvote'
											},
											id: chatId
										},
										{
											onSuccess: data => {
												if (data.vote) {
													successNotification({
														message: 'Voted!'
													})
												}
											}
										}
									)
									.catch(error => {
										console.error(error)
									})
							}}
						>
							<ThumbDownIcon />
						</Button>
					</TooltipTrigger>
					<TooltipContent>Downvote Response</TooltipContent>
				</Tooltip>
			</div>
		</TooltipProvider>
	)
}

export { MessageActions }
