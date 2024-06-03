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
import { signIn } from 'next-auth/react'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { z } from 'zod'

const formSchema = z.object({
	email: z.string().email({ message: 'Enter a valid email address' }),
	password: z.string().min(6, { message: 'Password must be at least 6 characters' })
})

type UserFormValue = z.infer<typeof formSchema>

export default function UserAuthForm() {
	const [loading] = useState(false)
	const defaultValues = {
		email: 'demo@gmail.com'
	}
	const form = useForm<UserFormValue>({
		resolver: zodResolver(formSchema),
		defaultValues
	})

	const onSubmit = async (data: UserFormValue) => {
		signIn('credentials', {
			email: data.email,
			callbackUrl: 'dashboard'
		})
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
						name="email"
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