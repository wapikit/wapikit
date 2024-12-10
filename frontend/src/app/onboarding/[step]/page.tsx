import { OnboardingSteps } from '~/constants'
import OnboardingStepClientPage from './client-page'
import { notFound } from 'next/navigation'

export function generateStaticParams() {
	return OnboardingSteps.map(step => ({
		step: step.slug
	}))
}

const OnboardingStepPage = async (props: any) => {
	const params = await props.params
	const stepSlug = params.step as string

	if (!OnboardingSteps.find(step => step.slug === stepSlug)) {
		notFound()
	}

	return <OnboardingStepClientPage stepSlug={stepSlug} />
}

export default OnboardingStepPage
