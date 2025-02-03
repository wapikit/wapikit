'use client'

import { useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { useAcceptOrganizationInvite, useGetOrganizationInviteBySlug } from 'root/.generated'
import { useLocalStorage } from '~/hooks/use-local-storage'
import { AUTH_TOKEN_LS } from '~/constants'
import { errorNotification } from '~/reusable-functions'
import { useAuthState } from '~/hooks/use-auth-state'
import { Button } from '~/components/ui/button'

const InvitationPage = () => {
	const router = useRouter()
	const [loading, setLoading] = useState(false)
	const setAuthToken = useLocalStorage<string | undefined>(AUTH_TOKEN_LS, undefined)[1]
	const { authState } = useAuthState()

	const params = useSearchParams()
	const slug = params.get('slug')

	const acceptOrganization = useAcceptOrganizationInvite()

	const { data: inviteData } = useGetOrganizationInviteBySlug(slug || '', {
		query: {
			enabled: !!authState.isAuthenticated
		}
	})

	const handleAccept = async () => {
		if (!slug) {
			router.push('/')
		}

		setLoading(true)
		try {
			const response = await acceptOrganization.mutateAsync({
				slug: slug || ''
			})

			if (response.token) {
				setAuthToken(response.token)
			}
			router.push('/dashboard')
		} catch (error) {
			errorNotification({
				message: 'Failed to accept invitation. Please try again.'
			})
		} finally {
			setLoading(false)
		}
	}

	return (
		<div className="flex min-h-screen flex-col items-center justify-center bg-gray-100 p-6">
			<div className="max-w-md rounded-lg bg-white p-6 text-center shadow-md">
				<h1 className="text-2xl font-bold text-gray-900">You're Invited!</h1>
				<p className="mt-4 text-gray-700">
					You have been invited to join{' '}
					<span className="font-semibold">{inviteData?.invite.organizationName}</span>.
				</p>
				<p className="mt-2 text-gray-600">Would you like to accept the invitation?</p>
				<div className="mt-6 flex justify-center gap-4">
					<Button onClick={handleAccept} disabled={loading}>
						{loading ? 'Processing...' : 'Accept'}
					</Button>
				</div>
			</div>
		</div>
	)
}

export default InvitationPage
