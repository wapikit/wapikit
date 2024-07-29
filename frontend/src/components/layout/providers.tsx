'use client'

import React from 'react'
import ThemeProvider from './ThemeToggle/theme-provider'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'

export default function Providers({ children }: { children: React.ReactNode }) {
	const queryClient = new QueryClient({
		defaultOptions: {
			mutations: {
				retry: false,
				onError(error, variables, context) {
					// if error is unauth access
					console.log({ error, variables, context })
				},
				throwOnError: false
			},
			queries: {
				retry: false,
				throwOnError: false,
				networkMode: 'online'
			}
		}
	})

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
