import { motion } from 'framer-motion'
import Link from 'next/link'

export const Overview = () => {
	return (
		<motion.div
			key="overview"
			className="mx-auto max-w-3xl md:mt-20"
			initial={{ opacity: 0, scale: 0.98 }}
			animate={{ opacity: 1, scale: 1 }}
			exit={{ opacity: 0, scale: 0.98 }}
			transition={{ delay: 0.5 }}
		>
			<div className="my-auto flex h-full max-w-xl flex-col justify-center gap-8 rounded-xl p-6 text-center leading-relaxed">
				<p>
					This is your{' '}
					<Link
						className="font-medium underline underline-offset-4"
						href="https://docs.wapikit.com/ai-chatbot"
						target="_blank"
					>
						organizational AI assistant.
					</Link>{' '}
					<br />
					You can ask your marketing campaign insights, conversational insights and much
					more.
				</p>
			</div>
		</motion.div>
	)
}
