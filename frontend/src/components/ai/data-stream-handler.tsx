'use client'

import { useChat } from 'ai/react'
import { useEffect, useRef } from 'react'
import { useAiChatStore } from '~/store/ai-chat-store'

type DataStreamDelta = {
	type:
		| 'text-delta'
		| 'code-delta'
		| 'title'
		| 'id'
		| 'suggestion'
		| 'clear'
		| 'finish'
		| 'user-message-id'
		| 'kind'
	content: string
}

export function DataStreamHandler({ id }: { id: string }) {
	const { data: dataStream , } = useChat({ id })
	const { chats } = useAiChatStore()
	const lastProcessedIndex = useRef(-1)

	useEffect(() => {
		if (!dataStream?.length) return
		console.log({ dataStream })

		const newDeltas = dataStream.slice(lastProcessedIndex.current + 1)
		lastProcessedIndex.current = dataStream.length - 1
		;(newDeltas as DataStreamDelta[]).forEach((delta: DataStreamDelta) => {
			if (delta.type === 'user-message-id') {
				return
			}




			// setBlock(draftBlock => {
			// 	if (!draftBlock) {
			// 		return { ...initialBlockData, status: 'streaming' }
			// 	}

			// 	switch (delta.type) {
			// 		case 'id':
			// 			return {
			// 				...draftBlock,
			// 				documentId: delta.content as string,
			// 				status: 'streaming'
			// 			}

			// 		case 'title':
			// 			return {
			// 				...draftBlock,
			// 				title: delta.content as string,
			// 				status: 'streaming'
			// 			}

			// 		case 'kind':
			// 			return {
			// 				...draftBlock,
			// 				kind: delta.content,
			// 				status: 'streaming'
			// 			}

			// 		case 'text-delta':
			// 			return {
			// 				...draftBlock,
			// 				content: draftBlock.content + (delta.content as string),
			// 				isVisible:
			// 					draftBlock.status === 'streaming' &&
			// 					draftBlock.content.length > 400 &&
			// 					draftBlock.content.length < 450
			// 						? true
			// 						: draftBlock.isVisible,
			// 				status: 'streaming'
			// 			}

			// 		case 'code-delta':
			// 			return {
			// 				...draftBlock,
			// 				content: delta.content as string,
			// 				isVisible:
			// 					draftBlock.status === 'streaming' &&
			// 					draftBlock.content.length > 300 &&
			// 					draftBlock.content.length < 310
			// 						? true
			// 						: draftBlock.isVisible,
			// 				status: 'streaming'
			// 			}

			// 		case 'clear':
			// 			return {
			// 				...draftBlock,
			// 				content: '',
			// 				status: 'streaming'
			// 			}

			// 		case 'finish':
			// 			return {
			// 				...draftBlock,
			// 				status: 'idle'
			// 			}

			// 		default:
			// 			return draftBlock
			// 	}
			// })
		})
	}, [dataStream])

	return null
}
