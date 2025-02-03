'use client'

import React from 'react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { errorNotification } from '~/reusable-functions'

export default function ApiQueryClientProvider({ children }: { children: React.ReactNode }) {
	const queryClient = new QueryClient({
		defaultOptions: {
			mutations: {
				retry: false,
				networkMode: 'online',
				onError: error => {
					if (((error as unknown as { status: number }).status as number) === 429) {
						errorNotification({
							message:
								'You have hit the rate limit. Please try again after some time.'
						})
					}
				}
			},
			queries: {
				retry: false,
				throwOnError: true,
				networkMode: 'online',
				refetchOnWindowFocus: false
			}
		}
	})

	return (
		<>
			<QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
		</>
	)
}
