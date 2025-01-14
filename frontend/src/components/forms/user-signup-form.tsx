'use client'
import { Button } from '~/components/ui/button'
import {
	Form,
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage
} from '~/components/ui/form'
import { Input } from '~/components/ui/input'
import { zodResolver } from '@hookform/resolvers/zod'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { useRegister, useVerifyOtp } from '~/generated'
import { useLocalStorage } from '~/hooks/use-local-storage'
import { AUTH_TOKEN_LS } from '~/constants'
import { errorNotification } from '~/reusable-functions'

const otpFormSchema = z.object({
	otp: z.string().length(6, { message: 'OTP must be 6 characters' })
})

const SignUpFormSchema = z.object({
	email: z.string().email({ message: 'Enter a valid email address' }),
	password: z.string().min(6, { message: 'Password must be at least 6 characters' }),
	confirmPassword: z.string().min(6, { message: 'Password must be at least 6 characters' }),
	name: z.string().min(6, { message: 'Name must be at least 6 characters' }),
	username: z.string().min(6, { message: 'Username must be at least 6 characters' }),
	orgInviteSlug: z.string().optional()
})

type SingUpFormValue = z.infer<typeof SignUpFormSchema>
type OtpFormValue = z.infer<typeof otpFormSchema>

export default function UserSignupForm() {
	const setAuthToken = useLocalStorage<string | undefined>(AUTH_TOKEN_LS, undefined)[1]

	const [isBusy, setIsBusy] = useState(false)
	const [activeForm, setActiveForm] = useState<'registrationDetailsForm' | 'otpForm'>(
		'registrationDetailsForm'
	)

	const defaultValues = {
		email: '',
		password: '',
		confirmPassword: '',
		name: '',
		orgInviteSlug: ''
	}

	const signUpForm = useForm<SingUpFormValue>({
		resolver: zodResolver(SignUpFormSchema),
		defaultValues
	})

	const otpForm = useForm<OtpFormValue>({
		resolver: zodResolver(otpFormSchema),
		defaultValues: {
			otp: ''
		}
	})

	const sendEmailConfirmationOtpMutation = useRegister()
	const createAccountMutation = useVerifyOtp()

	async function initiateRegistration(data: SingUpFormValue) {
		try {
			if (isBusy) {
				return
			}
			setIsBusy(true)

			if (data.password !== data.confirmPassword) {
				errorNotification({
					message: 'Passwords do not match'
				})
				return
			}

			const response = await sendEmailConfirmationOtpMutation.mutateAsync({
				data: {
					password: data.password,
					username: data.email,
					email: data.email,
					name: data.name,
					organizationInviteSlug: data.orgInviteSlug || undefined
				}
			})

			if (response.isOtpSent) {
				// open the otp form
				setActiveForm(() => 'otpForm')
			} else {
				// something went wrong show error token not found
			}
		} catch (error) {
			console.error(error)
			errorNotification({
				message: 'Something went wrong while creating your account'
			})
		} finally {
			setIsBusy(false)
		}
	}

	async function submitOtp(data: OtpFormValue) {
		try {
			if (isBusy) {
				return
			}
			setIsBusy(true)

			const userData = signUpForm.getValues()

			const response = await createAccountMutation.mutateAsync({
				data: {
					password: userData.password,
					username: userData.username,
					email: userData.email,
					name: userData.name,
					organizationInviteSlug: userData.orgInviteSlug || undefined,
					otp: data.otp
				}
			})

			if (response.token) {
				setAuthToken(response.token)
				window.location.href = '/dashboard'
			} else {
				// something went wrong show error token not found
				errorNotification({
					message: 'Something went wrong while creating your account'
				})
			}
		} catch (error) {
			console.error(error)
			errorNotification({
				message: 'Something went wrong while creating your account'
			})
		} finally {
			setIsBusy(false)
		}
	}

	return (
		<>
			{activeForm === 'registrationDetailsForm' ? (
				<Form {...signUpForm}>
					<form
						onSubmit={signUpForm.handleSubmit(initiateRegistration)}
						className="flex w-full flex-col gap-2 space-y-2"
						id="registration-details-form"
					>
						<FormField
							control={signUpForm.control}
							name="email"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Email</FormLabel>
									<FormControl>
										<Input
											type="email"
											placeholder="Enter your email..."
											disabled={isBusy}
											{...field}
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={signUpForm.control}
							name="name"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Name</FormLabel>
									<FormControl>
										<Input
											placeholder="Enter your name"
											disabled={isBusy}
											{...field}
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={signUpForm.control}
							name="username"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Username</FormLabel>
									<FormControl>
										<Input
											placeholder="Enter your username"
											disabled={isBusy}
											{...field}
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={signUpForm.control}
							name="password"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Password</FormLabel>
									<FormControl>
										<Input
											type="password"
											placeholder="Enter your password..."
											disabled={isBusy}
											{...field}
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={signUpForm.control}
							name="confirmPassword"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Confirm Password</FormLabel>
									<FormControl>
										<Input
											type="password"
											placeholder="Confirm your password"
											disabled={isBusy}
											{...field}
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={signUpForm.control}
							name="orgInviteSlug"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Organization Invitation Id (Optional)</FormLabel>
									<FormControl>
										<Input
											placeholder="#########"
											disabled={isBusy}
											{...field}
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						<Button disabled={isBusy} className="ml-auto w-full" type="submit">
							Confirm Email
						</Button>
					</form>
				</Form>
			) : (
				<Form {...otpForm}>
					<form
						onSubmit={otpForm.handleSubmit(submitOtp)}
						className="flex w-full flex-col gap-2 space-y-2"
						id="otp-form"
					>
						<FormField
							control={otpForm.control}
							name="otp"
							key={'otp_field'}
							render={({ field }) => (
								<FormItem>
									<FormLabel>Otp</FormLabel>
									<FormControl>
										<Input
											id="otp"
											placeholder="Enter OTP"
											disabled={isBusy}
											{...field}
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						<Button disabled={isBusy} className="ml-auto w-full" type="submit">
							Sign Up
						</Button>
					</form>
				</Form>
			)}
		</>
	)
}
