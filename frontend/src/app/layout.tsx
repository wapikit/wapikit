import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'
import Providers from '~/components/layout/providers'
import { Toaster } from '~/components/ui/toaster'
import NextTopLoader from 'nextjs-toploader'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
	// ! TODO: fetch this data from backend only
	title: 'WapiKit'
}

export default async function RootLayout({
	children
}: Readonly<{
	children: React.ReactNode
}>) {
	const session = null
	return (
		<html lang="en">
			<body className={inter.className}>
				<NextTopLoader />
				<Providers session={session}>
					<Toaster />
					{children}
				</Providers>
			</body>
		</html>
	)
}
