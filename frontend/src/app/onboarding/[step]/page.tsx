import { OnboardingSteps } from '~/constants'
import OnboardingStepClientPage from './client-page'

export function generateStaticParams() {
	return OnboardingSteps.map(step => ({
		step: step.slug
	}))
}

const OnboardingStepPage = async (props: any) => {
	const params = await props.params
	const stepSlug = params.step as string

	return <OnboardingStepClientPage stepSlug={stepSlug} />
}

export default OnboardingStepPage
