'use client'

import { clsx } from 'clsx'
import { CheckIcon } from 'lucide-react'
import Image from 'next/image'
import Link from 'next/link'
import { useParams, useRouter } from 'next/navigation'
import { Icons } from '~/components/icons'
import { buttonVariants } from '~/components/ui/button'
import { OnboardingSteps } from '~/constants'
import { useLayoutStore } from '~/store/layout.store'

const OnboardingLayout = (props: { children: React.ReactNode }) => {
	const { onboardingSteps } = useLayoutStore()

	console.log('onboardingSteps', onboardingSteps)

	const params = useParams()
	const step = params.step as string
	const currentStep = onboardingSteps.find(s => s.slug === step)

	if (!currentStep) {
		return <>{props.children}</>
	}

	return (
		<div className="relative h-screen flex-col items-center justify-center md:grid lg:max-w-none lg:grid-cols-2 lg:px-0">
			<Link
				href="/signin"
				className={clsx(
					buttonVariants({ variant: 'ghost' }),
					'absolute right-4 top-4 hidden md:right-8 md:top-8'
				)}
			>
				Login
			</Link>
			<div className="relative hidden h-full flex-col bg-muted p-10 text-white dark:border-r lg:flex">
				<div className="absolute inset-0 bg-gradient-to-br from-[rgb(2,105,67)] via-[rgb(0,3,2)] to-[rgb(28,68,72)]" />
				<div className="relative z-20 flex items-center text-lg font-medium">
					<Image src={'/logo/dark.svg'} width={100} height={40} alt="logo" />
				</div>

				{/* ==== onboarding steps UI ===== */}
				<ol role="list" className="my-auto pt-28">
					{onboardingSteps.map((step, stepIdx) => {
						const IconToRender = Icons[step.icon]

						console.log('step', step)

						return (
							<li
								key={step.slug}
								className={clsx(
									stepIdx !== OnboardingSteps.length - 1 ? 'pb-16' : '',
									'relative'
								)}
							>
								{step.status === 'complete' ? (
									<>
										{stepIdx !== OnboardingSteps.length - 1 ? (
											<div
												className="absolute left-4 top-4 -ml-px mt-0.5 h-full w-0.5 bg-[#25d366]"
												aria-hidden="true"
											/>
										) : null}
										<div className="group relative flex items-start">
											<span className="flex h-9 items-center">
												<span className="relative z-10 flex items-center justify-center rounded-md border border-[#25d366] p-2 backdrop-blur-md">
													<CheckIcon
														className={`size-5 text-[#25d366]`}
													/>
												</span>
											</span>
											<span className="ml-4 flex min-w-0 flex-col">
												<span className="text-xs font-semibold uppercase tracking-wide text-[#25d366]">
													{step.title}
												</span>
												<span className="text-sm opacity-55">
													{step.description}
												</span>
											</span>
										</div>
									</>
								) : step.status === 'current' ? (
									<>
										{stepIdx !== OnboardingSteps.length - 1 ? (
											<div
												className="absolute left-4 top-4 -ml-px mt-0.5 h-full w-0.5 bg-gray-300"
												aria-hidden="true"
											/>
										) : null}
										<div
											className="group relative flex items-start"
											aria-current="step"
										>
											<span className="flex h-9 items-center">
												<span className="relative z-10 flex items-center justify-center rounded-md border border-primary-foreground p-2 backdrop-blur-md">
													<IconToRender className={`size-5`} />
												</span>
											</span>
											<span className="ml-4 flex min-w-0 flex-col">
												<span className="text-xs font-semibold uppercase tracking-wide text-primary-foreground">
													{step.title}
												</span>
												<span className="text-sm opacity-55">
													{step.description}
												</span>
											</span>
										</div>
									</>
								) : (
									<>
										{stepIdx !== OnboardingSteps.length - 1 ? (
											<div
												className="absolute left-4 top-4 -ml-px mt-0.5 h-full w-0.5 bg-gray-300"
												aria-hidden="true"
											/>
										) : null}
										<div className="group relative flex items-start">
											<span className="flex h-9 items-center">
												<span className="relative z-10 flex items-center justify-center rounded-md border border-primary-foreground p-2 backdrop-blur-lg">
													<IconToRender className={`size-5`} />
												</span>
											</span>
											<span className="ml-4 flex min-w-0 flex-col">
												<span className="text-xs font-semibold uppercase tracking-wide ">
													{step.title}
												</span>
												<span className="text-sm opacity-55">
													{step.description}
												</span>
											</span>
										</div>
									</>
								)}
							</li>
						)
					})}
				</ol>

				<div className="relative z-20 mt-auto">
					<blockquote className="space-y-2">
						<p className="text-sm">
							⭐️ Star us on github{' '}
							<Link
								href={'https://github.com/sarthakjdev/wapikit'}
								className="underline"
								target="_blank"
							>
								here
							</Link>
							.
						</p>
					</blockquote>
				</div>
			</div>
			<div className="flex h-full items-center p-4 lg:p-8">
				<div className="mx-auto flex w-full max-w-xl flex-col justify-center space-y-6">
					<div className="flex flex-col space-y-2 text-left">
						<h1 className="text-xl font-semibold tracking-tight">
							{currentStep.title}
						</h1>
						<p className="text-sm">{currentStep.description}</p>
					</div>
					{props.children}
				</div>
			</div>
		</div>
	)
}

export default OnboardingLayout
