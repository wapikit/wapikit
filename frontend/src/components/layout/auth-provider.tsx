'use client'

import { usePathname, useRouter } from 'next/navigation'
import React, { useEffect, useRef } from 'react'
import { useAuthState } from '~/hooks/use-auth-state'
import LoadingSpinner from '../loader'
import { useGetAllPhoneNumbers, useGetAllTemplates, useGetUser } from 'root/.generated'
import { useLayoutStore } from '~/store/layout.store'
import { OnboardingStepsEnum } from '~/constants'
import CreateTagModal from '../forms/create-tag'

const AuthProvisioner: React.FC<{ children: React.ReactNode }> = ({ children }) => {
	const { authState } = useAuthState()
	const router = useRouter()
	const pathname = usePathname()

	const { writeProperty, onboardingSteps, currentOrganization, user } = useLayoutStore()

	useEffect(() => {
		if (pathname === '/signin' || pathname === '/logout' || pathname === '/signup') {
			return
		} else {
			if (authState.isAuthenticated === false) {
				router.push('/signin')
			} else {
				// either auth is loading or user is authenticated
				if (pathname === '/') {
					router.push('/dashboard')
				}
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
		if (phoneNumbersResponse && templatesResponse) {
			writeProperty({
				phoneNumbers: phoneNumbersResponse,
				templates: templatesResponse
			})
		}
	}, [phoneNumbersResponse, templatesResponse, writeProperty])

	const isRedirecting = useRef(false)

	useEffect(() => {
		if (isRedirecting.current) return

		if (!authState.isAuthenticated || !userData || (user && currentOrganization)) {
			return
		}

		if (!authState.data.user.organizationId) {
			writeProperty({
				user: userData.user,
				onboardingSteps: onboardingSteps.map(step => {
					if (step.slug === OnboardingStepsEnum.CreateOrganization) {
						return {
							...step,
							status: 'current'
						}
					} else {
						return step
					}
				})
			})

			router.push(`/onboarding/${OnboardingStepsEnum.CreateOrganization}`)
			isRedirecting.current = true
		} else {
			writeProperty({
				user: userData.user,
				currentOrganization: userData.user.organization,
				isOwner: userData.user.currentOrganizationAccessLevel === 'Owner'
			})

			if (!userData.user.organization?.whatsappBusinessAccountDetails) {
				writeProperty({
					onboardingSteps: onboardingSteps.map(step => {
						if (step.slug === OnboardingStepsEnum.CreateOrganization) {
							return {
								...step,
								status: 'complete'
							}
						} else if (
							step.slug === OnboardingStepsEnum.WhatsappBusinessAccountDetails
						) {
							return {
								...step,
								status: 'current'
							}
						} else {
							return step
						}
					})
				})

				router.push(`/onboarding/${OnboardingStepsEnum.WhatsappBusinessAccountDetails}`)
				isRedirecting.current = true
			}
		}
	}, [userData, writeProperty, authState, onboardingSteps, router, user, currentOrganization])

	if (
		typeof authState.isAuthenticated !== 'boolean' &&
		!authState.isAuthenticated &&
		pathname !== '/'
	) {
		return (
			<div className="flex h-full w-full items-center justify-center">
				<LoadingSpinner />
				<CreateTagModal
					setIsCreateTagModalOpen={value => {
						writeProperty({
							isCreateTagModalOpen: value
						})
					}}
				/>
			</div>
		)
	} else {
		return <>{children}</>
	}
}

export default AuthProvisioner
