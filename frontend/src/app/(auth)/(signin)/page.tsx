import { type Metadata } from 'next'
import Link from 'next/link'
import UserAuthForm from '~/components/forms/user-auth-form'
import { buttonVariants } from '~/components/ui/button'
import { clsx } from 'clsx'
import Image from 'next/image'

export const metadata: Metadata = {
	title: 'Signin | WapiKit'
}

export default function AuthenticationPage() {
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
				<div className="absolute inset-0 bg-zinc-900" />
				<div className="relative z-20 flex items-center text-lg font-medium">
					<Image src={'/logo/dark.svg'} width={100} height={40} alt="logo" />
				</div>
				<div className="relative z-20 mt-auto">
					<blockquote className="space-y-2">
						<p className="text-lg">
							Star us on github. You can report bugs on github issues{' '}
							<Link href={'https://github.com/sarthakjdev/wapikit'}>here</Link>.
						</p>
					</blockquote>
				</div>
			</div>
			<div className="flex h-full items-center p-4 lg:p-8">
				<div className="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
					<div className="flex flex-col space-y-2 text-center">
						<h1 className="text-2xl font-semibold tracking-tight">
							Sign in to WapiKit
						</h1>
						<p className="text-sm text-muted-foreground">
							Enter your email below to create your account
						</p>
					</div>
					<UserAuthForm />
					<p className="px-8 text-center text-sm text-muted-foreground">
						By clicking continue, you agree to our{' '}
						<Link
							href="/terms"
							className="underline underline-offset-4 hover:text-primary"
						>
							Terms of Service
						</Link>{' '}
						and{' '}
						<Link
							href="/privacy"
							className="underline underline-offset-4 hover:text-primary"
						>
							Privacy Policy
						</Link>
						.
					</p>
				</div>
			</div>
		</div>
	)
}
