'use client'

import { usePathname, useRouter } from 'next/navigation'
import React, { useEffect } from 'react'
import { useAuthState } from '~/hooks/use-auth-state'
import LoadingSpinner from '../loader'
import { useGetAllPhoneNumbers, useGetAllTemplates, useGetUser } from 'root/.generated'
import { useLayoutStore } from '~/store/layout.store'

const AuthProvisioner: React.FC<{ children: React.ReactNode }> = ({ children }) => {
	const { authState } = useAuthState()
	const router = useRouter()
	const pathname = usePathname()

	const { writeProperty, onboardingSteps } = useLayoutStore()

	useEffect(() => {
		if (pathname === '/signin') {
			return
		} else {
			if (authState.isAuthenticated === false) {
				router.push('/signin')
			} else {
				// either auth is loading or user is authenticated
			}
		}
	}, [authState.isAuthenticated, pathname, router])

	const { data: userData } = useGetUser({
		query: {
			enabled: !!authState.isAuthenticated
		}
	})

	const { data: phoneNumbersResponse } = useGetAllPhoneNumbers({
		query: {
			enabled: !!authState.isAuthenticated
		}
	})

	const { data: templatesResponse } = useGetAllTemplates({
		query: {
			enabled: !!authState.isAuthenticated
		}
	})

	useEffect(() => {
		if (!authState.isAuthenticated || !userData) {
			return
		}

		if (!authState.data.user.organizationId) {
			writeProperty({
				user: userData.user,
				onboardingSteps: onboardingSteps.map(step => {
					if (step.slug === 'create-organization') {
						return {
							...step,
							status: 'current'
						}
					} else {
						return step
					}
				})
			})

			router.push('/onboarding/create-organization')
		} else {
			writeProperty({
				user: userData.user,
				currentOrganization: userData.user.organization,
				isOwner: userData.user.currentOrganizationAccessLevel === 'Owner'
			})

			if (!userData.user.organization?.whatsappBusinessAccountDetails) {
				writeProperty({
					onboardingSteps: onboardingSteps.map(step => {
						if (step.slug === 'create-organization') {
							return {
								...step,
								status: 'complete'
							}
						} else if (step.slug === 'whatsapp-business-account-details') {
							return {
								...step,
								status: 'current'
							}
						} else {
							return step
						}
					})
				})

				router.push(`/onboarding/whatsapp-business-account-details`)
			}
		}

		if (phoneNumbersResponse && templatesResponse) {
			writeProperty({
				phoneNumbers: phoneNumbersResponse,
				templates: templatesResponse
			})
		}
	}, [
		userData,
		authState.isAuthenticated,
		writeProperty,
		phoneNumbersResponse,
		templatesResponse
	])

	if (
		typeof authState.isAuthenticated !== 'boolean' &&
		!authState.isAuthenticated &&
		pathname !== '/'
	) {
		return (
			<div className="flex h-full w-full items-center justify-center">
				<LoadingSpinner />
			</div>
		)
	} else {
		return <>{children}</>
	}
}

export default AuthProvisioner
