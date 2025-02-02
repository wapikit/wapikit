'use client'

import { Inter } from 'next/font/google'
import './globals.css'
import ApiQueryClientProvider from '~/components/layout/api-query-client-provider'
import { Toaster } from '~/components/ui/sonner'
import NextTopLoader from 'nextjs-toploader'
import AuthProvisioner from '~/components/layout/auth-provider'
import SseConnectionProvider from '~/components/layout/sse-provider'
import MetaProvider from '~/components/layout/meta-provider'
import AiChatProvider from '~/components/layout/ai-chat-provider'
import dynamic from 'next/dynamic'
import CommandMenuProvider from '~/components/layout/command-menu-provider'
import { clsx } from 'clsx'

const inter = Inter({ subsets: ['latin'] })

const ThemeProvider = dynamic(() => import('../components/layout/theme/theme-provider'), {
	ssr: false
})

export default function RootLayout({
	children
}: Readonly<{
	children: React.ReactNode
}>) {
	return (
		<html lang="en">
			<body className={clsx(inter.className, '!font-sans antialiased')}>
				<NextTopLoader />
				<ThemeProvider attribute="class" defaultTheme="system" enableSystem>
					<ApiQueryClientProvider>
						<Toaster />
						<AuthProvisioner>
							<MetaProvider>
								<SseConnectionProvider>
									<AiChatProvider>
										<CommandMenuProvider />
										{children}
									</AiChatProvider>
								</SseConnectionProvider>
							</MetaProvider>
						</AuthProvisioner>
					</ApiQueryClientProvider>
				</ThemeProvider>
				<div className="display: none;">
					<span className=" stroke-green-500"></span>
					<span className="stroke-red-500"></span>
					<span className="stroke-blue-500"></span>
					<span className="stroke-yellow-500"></span>
					<span className="fill-green-500"></span>
					<span className="fill-indigo-500 "></span>
					<span className="fill-red-500 "></span>
					{/* Add other colors as needed */}
				</div>
			</body>
		</html>
	)
}
