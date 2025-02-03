'use client'

import Link from 'next/link'
import UserLoginForm from '~/components/forms/user-login-form'
import { buttonVariants } from '~/components/ui/button'
import { clsx } from 'clsx'
import Image from 'next/image'
import { useAuthState } from '~/hooks/use-auth-state'
import { redirect } from 'next/navigation'
import LoadingSpinner from '~/components/loader'

export default function AuthenticationPage() {
	const { authState } = useAuthState()

	if (authState.isAuthenticated) {
		redirect('/dashboard')
	} else if (authState.isAuthenticated === false) {
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

					<div className="relative z-20 mt-auto text-left text-2xl font-bold leading-relaxed  md:text-3xl">
						<span className="text-4xl font-semibold">
							Manage your WhatsApp Business Comms.
						</span>{' '}
						<br />
						<span className="text-lg italic">AI-Driven.</span>
						<span className="text-lg italic"> Data Privacy.</span>
						<span className="text-lg italic"> 100% Control.</span>
					</div>

					<div className="relative z-20 mt-auto">
						<blockquote className="space-y-2">
							<p className="text-sm">
								⭐️ Star us on{' '}
								<Link
									href={'https://github.com/wapikit/wapikit'}
									className="underline"
									target="_blank"
								>
									Github
								</Link>
								.
							</p>
						</blockquote>
					</div>
				</div>
				<div className="flex h-full items-center p-4 lg:p-8">
					<div className="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
						<div className="flex flex-col space-y-2 text-left">
							<h1 className="text-2xl font-semibold tracking-tight">
								Sign in to WapiKit
							</h1>
						</div>
						<UserLoginForm />
						<p className="text-left text-sm text-muted-foreground">
							By clicking continue, you agree to our{' '}
							<Link
								href="/terms-of-service"
								className="underline underline-offset-4 hover:text-primary"
							>
								Terms of Service
							</Link>{' '}
							and{' '}
							<Link
								href="/privacy-policy"
								className="underline underline-offset-4 hover:text-primary"
							>
								Privacy Policy
							</Link>
							.
						</p>
						<p className="text-left text-xs text-muted-foreground">
							Don't have an account?{' '}
							<Link
								href="/signup"
								className="underline underline-offset-4 hover:text-primary"
							>
								Signup
							</Link>
							.
						</p>
					</div>
				</div>
			</div>
		)
	} else {
		// auth is still loading
		return <LoadingSpinner />
	}
}
