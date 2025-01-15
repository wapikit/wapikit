'use client'

import { memo } from 'react'

function PureChatHeader({ chatTitle }: { chatTitle: string }) {
	return (
		<header className="sticky top-0 flex items-center gap-2 bg-background px-2 py-1.5 md:px-2">
			{chatTitle}
		</header>
	)
}

export const ChatHeader = memo(PureChatHeader, () => {
	return true
})
