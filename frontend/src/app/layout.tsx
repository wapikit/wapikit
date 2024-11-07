import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'
import Providers from '~/components/layout/providers'
import { Toaster } from '~/components/ui/sonner'
import NextTopLoader from 'nextjs-toploader'
import AuthProvisioner from '~/components/layout/auth'
import WebsocketConnectionProvider from '~/components/layout/websocket'

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
	return (
		<html lang="en">
			<body className={inter.className}>
				<NextTopLoader />
				<Providers>
					<Toaster />
					<AuthProvisioner>
						<WebsocketConnectionProvider>{children}</WebsocketConnectionProvider>
					</AuthProvisioner>
				</Providers>
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
