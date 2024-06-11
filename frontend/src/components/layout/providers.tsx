'use client'
import React from 'react'
import ThemeProvider from './ThemeToggle/theme-provider'
import { SessionProvider, type SessionProviderProps } from 'next-auth/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'

export default function Providers({
	session,
	children
}: {
	session: SessionProviderProps['session']
	children: React.ReactNode
}) {
	const queryClient = new QueryClient()

	return (
		<>
			<QueryClientProvider client={queryClient}>
				<ThemeProvider attribute="class" defaultTheme="system" enableSystem>
					<SessionProvider session={session}>{children}</SessionProvider>
				</ThemeProvider>
			</QueryClientProvider>
		</>
	)
}
