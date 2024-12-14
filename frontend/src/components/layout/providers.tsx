'use client'

import React from 'react'
import dynamic from 'next/dynamic'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { errorNotification } from '~/reusable-functions'

const ThemeProvider = dynamic(() => import('./ThemeToggle/theme-provider'), { ssr: false })

export default function Providers({ children }: { children: React.ReactNode }) {
	const queryClient = new QueryClient({
		defaultOptions: {
			mutations: {
				retry: false,
				onError(error, variables, context) {
					// if error is unauth access

					// if errorCode === 402 . router.push('/logout)

					errorNotification({
						message: error.message
					})

					console.log({ error, variables, context })
				},
				throwOnError: true,
				networkMode: 'online'
			},
			queries: {
				retry: false,
				throwOnError: true,
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
