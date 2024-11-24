'use client'

import { useRouter } from 'next/navigation'
import { useEffect } from 'react'
import LoadingSpinner from '~/components/loader'
import { useLayoutStore } from '~/store/layout.store'

const OnboardingPage = () => {
	const { onboardingSteps } = useLayoutStore()
	const router = useRouter()

	useEffect(() => {
		console.log('onboardingSteps', onboardingSteps)

		const stepToRedirectTo = onboardingSteps.find(
			step => step.status === 'current' || step.status === 'incomplete'
		)

		console.log('stepToRedirectTo', stepToRedirectTo)

		if (stepToRedirectTo) {
			router.push(`/onboarding/${stepToRedirectTo.slug}`)
		} else {
			router.push('/dashboard')
		}
	}, [onboardingSteps])

	return <LoadingSpinner />
}

export default OnboardingPage
