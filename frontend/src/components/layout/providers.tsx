'use client'
import React from 'react'
import ThemeProvider from './ThemeToggle/theme-provider'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'

export default function Providers({ children }: { children: React.ReactNode }) {
	const queryClient = new QueryClient()

	return (
		<>
			<QueryClientProvider client={queryClient}>
				<ThemeProvider attribute="class" defaultTheme="system" enableSystem>
					{children}
				</ThemeProvider>
			</QueryClientProvider>
		</>
	)
}
