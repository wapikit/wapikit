import Header from '~/components/layout/header'
import Sidebar from '~/components/layout/sidebar'
import type { Metadata } from 'next'
import CreateTagModal from '~/components/forms/create-tag'

export const metadata: Metadata = {
	title: 'Wapikit',
	description: ''
}

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
	return (
		<>
			<Header />
			<CreateTagModal />
			<div className="flex h-screen overflow-hidden">
				<Sidebar />
				<main className="flex-1 overflow-hidden pt-16">{children}</main>
			</div>
		</>
	)
}
