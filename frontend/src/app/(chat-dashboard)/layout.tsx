import type { Metadata } from 'next'

export const metadata: Metadata = {
	title: 'Chat | WapiKit'
}

export default async function RootLayout({
	children
}: Readonly<{
	children: React.ReactNode
}>) {
	return <>{children}</>
}
