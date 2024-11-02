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
import { useLogin } from '~/generated'
import { useLocalStorage } from '~/hooks/use-local-storage'
import { AUTH_TOKEN_LS } from '~/constants'
import { useRouter } from 'next/navigation'

const formSchema = z.object({
	email: z.string().email({ message: 'Enter a valid email address' }),
	password: z.string().min(6, { message: 'Password must be at least 6 characters' })
})

type UserFormValue = z.infer<typeof formSchema>

export default function UserAuthForm() {
	const setAuthToken = useLocalStorage<string | undefined>(AUTH_TOKEN_LS, undefined)[1]

	const router = useRouter()

	const [loading] = useState(false)

	const defaultValues = {
		email: '',
		password: ''
	}
	const form = useForm<UserFormValue>({
		resolver: zodResolver(formSchema),
		defaultValues
	})

	const mutation = useLogin()

	const onSubmit = async (data: UserFormValue) => {
		await mutation.mutateAsync(
			{
				data: {
					password: data.password,
					username: data.email
				}
			},
			{
				onSuccess: data => {
					if (data.token) {
						setAuthToken(data.token)
						router.push('/dashboard')
					} else {
						// something went wrong show error token not found
					}
				},
				onError: error => {
					console.error(error)
				}
			}
		)
	}

	return (
		<>
			<Form {...form}>
				<form
					onSubmit={form.handleSubmit(onSubmit)}
					className="flex w-full flex-col gap-2 space-y-2"
				>
					<FormField
						control={form.control}
						name="email"
						render={({ field }) => (
							<FormItem>
								<FormLabel>Email</FormLabel>
								<FormControl>
									<Input
										type="email"
										placeholder="Enter your email..."
										disabled={loading}
										{...field}
									/>
								</FormControl>
								<FormMessage />
							</FormItem>
						)}
					/>

					<FormField
						control={form.control}
						name="password"
						render={({ field }) => (
							<FormItem>
								<FormLabel>Password</FormLabel>
								<FormControl>
									<Input
										type="password"
										placeholder="Enter your password..."
										disabled={loading}
										{...field}
									/>
								</FormControl>
								<FormMessage />
							</FormItem>
						)}
					/>

					<Button
						disabled={loading}
						onClick={() => {}}
						className="ml-auto w-full"
						type="submit"
					>
						Sign in
					</Button>
				</form>
			</Form>
		</>
	)
}
