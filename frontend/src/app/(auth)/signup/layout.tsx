import type { Metadata } from 'next'

export const metadata: Metadata = {
	title: 'Signin | WapiKit'
}

export default function SignupPageLayout({
	children
}: Readonly<{
	children: React.ReactNode
}>) {
	return <> {children}</>
}
