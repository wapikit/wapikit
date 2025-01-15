'use client'

import { motion } from 'framer-motion'
import { Button } from '~/components/ui/button'
import { memo } from 'react'
import { useAiChatStore } from '~/store/ai-chat-store'

interface SuggestedActionsProps {
	selectSuggestedAction: (action: string) => void
}

function PureSuggestedActions({ selectSuggestedAction }: SuggestedActionsProps) {
	const { suggestions } = useAiChatStore()

	return (
		<div className="grid w-full gap-2 sm:grid-cols-2">
			{suggestions.map((suggestedAction, index) => (
				<motion.div
					initial={{ opacity: 0, y: 20 }}
					animate={{ opacity: 1, y: 0 }}
					exit={{ opacity: 0, y: 20 }}
					transition={{ delay: 0.05 * index }}
					key={`suggested-action-${suggestedAction.title}-${index}`}
					className={index > 1 ? 'hidden sm:block' : 'block'}
				>
					<Button
						variant="ghost"
						onClick={async () => {
							selectSuggestedAction(suggestedAction.action)
						}}
						className="h-auto w-full flex-1 items-start justify-start gap-1 rounded-md border px-4 py-3.5 text-left text-sm sm:flex-col"
					>
						<span className="font-medium">{suggestedAction.title}</span>
						<span className="text-muted-foreground">{suggestedAction.label}</span>
					</Button>
				</motion.div>
			))}
		</div>
	)
}

export const SuggestedActions = memo(PureSuggestedActions, () => true)
