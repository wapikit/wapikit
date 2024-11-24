import { OnboardingSteps } from '~/constants'
import OnboardingStepClientPage from './client-page'

export function generateStaticParams() {
	return OnboardingSteps.map(step => ({
		step: step.slug
	}))
}

const OnboardingStepPage = ({ params }: { params: { step: string } }) => {
	const stepSlug = params.step

	return <OnboardingStepClientPage stepSlug={stepSlug} />
}

export default OnboardingStepPage
