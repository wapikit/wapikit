'use client'

import { useRouter } from 'next/navigation'
import { useEffect } from 'react'
import LoadingSpinner from '~/components/loader'
import { useLayoutStore } from '~/store/layout.store'

const OnboardingPage = () => {
	const { onboardingSteps } = useLayoutStore()
	const router = useRouter()

	useEffect(() => {
		const stepToRedirectTo = onboardingSteps.find(
			step => step.status === 'current' || step.status === 'incomplete'
		)

		if (stepToRedirectTo) {
			router.push(`/onboarding/${stepToRedirectTo.slug}`)
		} else {
			router.push('/dashboard')
		}
	}, [onboardingSteps, router])

	return <LoadingSpinner />
}

export default OnboardingPage
