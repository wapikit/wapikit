'use client'

import { zodResolver } from '@hookform/resolvers/zod'
import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import {
	useCreateOrganization,
	useCreateOrganizationInvite,
	UserPermissionLevel,
	useUpdateWhatsappBusinessAccountDetails,
	switchOrganization
} from 'root/.generated'
import { type z } from 'zod'
import { errorNotification, materialConfirm, successNotification } from '~/reusable-functions'
import {
	NewOrganizationFormSchema,
	NewTeamMemberInviteFormSchema,
	WhatsappBusinessAccountDetailsFormSchema
} from '~/schema'
import { useLayoutStore } from '~/store/layout.store'
import { EyeIcon } from 'lucide-react'
import {
	Form,
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage
} from '~/components/ui/form'
import { Input } from '~/components/ui/input'
import { Button } from '~/components/ui/button'
import { Textarea } from '~/components/ui/textarea'
import { useRouter } from 'next/navigation'
import { AUTH_TOKEN_LS } from '~/constants'

const OnboardingStepClientPage = ({ stepSlug }: { stepSlug: string }) => {
	const router = useRouter()

	const createOrganizationMutation = useCreateOrganization()
	const inviteUserMutation = useCreateOrganizationInvite()
	const updateWhatsappBusinessAccountDetailsMutation = useUpdateWhatsappBusinessAccountDetails()

	const { currentOrganization, onboardingSteps, writeProperty } = useLayoutStore()

	const newMemberInviteForm = useForm<z.infer<typeof NewTeamMemberInviteFormSchema>>({
		resolver: zodResolver(NewTeamMemberInviteFormSchema),
		defaultValues: {
			email: ''
		}
	})

	const newOrganizationForm = useForm<z.infer<typeof NewOrganizationFormSchema>>({
		resolver: zodResolver(NewOrganizationFormSchema),
		defaultValues: {
			name: '',
			description: ''
		}
	})

	const whatsappBusinessAccountForm = useForm<
		z.infer<typeof WhatsappBusinessAccountDetailsFormSchema>
	>({
		resolver: zodResolver(WhatsappBusinessAccountDetailsFormSchema),
		defaultValues: {
			whatsappBusinessAccountId: currentOrganization?.businessAccountId || undefined,
			apiToken: currentOrganization?.whatsappBusinessAccountDetails?.accessToken || undefined,
			webhookSecret:
				currentOrganization?.whatsappBusinessAccountDetails?.webhookSecret || undefined
		}
	})

	const [isBusy, setIsBusy] = useState(false)
	const [whatsAppBusinessAccountDetailsVisibility, setWhatsAppBusinessAccountDetailsVisibility] =
		useState({
			whatsappBusinessAccountId: false,
			apiToken: false,
			webhookSecret: false
		})

	useEffect(() => {
		const step = onboardingSteps.find(step => step.slug === stepSlug)

		if (step?.status === 'complete') {
			const nextStep = onboardingSteps.find(
				(_, index) => index === onboardingSteps.findIndex(s => s.slug === stepSlug) + 1
			)

			if (nextStep) {
				router.push(`/onboarding/${nextStep.slug}`)
			} else {
				router.push('/dashboard')
			}
		}
	}, [onboardingSteps, router, stepSlug])

	const step = onboardingSteps.find(step => step.slug === stepSlug)

	if (!step) {
		return <div>Step not found</div>
	}

	async function inviteTeamMembers() {
		try {
			setIsBusy(true)
			const confirmation = await materialConfirm({
				description: 'Are you sure you want to invite this user?',
				title: 'Invite User'
			})

			if (!confirmation) return

			const response = await inviteUserMutation.mutateAsync({
				data: {
					email: newMemberInviteForm.getValues('email'),
					accessLevel: UserPermissionLevel.Member
				}
			})

			if (response.invite) {
				successNotification({
					message: 'User invited successfully.'
				})
				newMemberInviteForm.reset()
			} else {
				errorNotification({
					message: 'Something went wrong, While inviting a user. Please try again.'
				})
			}
		} catch (error) {
			console.error(error)
			errorNotification({
				message: 'Something went wrong, While inviting a user. Please try again.'
			})
		} finally {
			setIsBusy(false)
		}
	}

	async function handleCreateOrganization(data: z.infer<typeof NewOrganizationFormSchema>) {
		try {
			const response = await createOrganizationMutation.mutateAsync({
				data: {
					name: data.name
					// description: data.description || undefined
				}
			})

			if (response.organization) {
				const switched = await switchOrganization({
					organizationId: response.organization.uniqueId
				})

				if (switched.token) {
					writeProperty({
						currentOrganization: response.organization
					})

					window.localStorage.setItem(AUTH_TOKEN_LS, switched.token)
					window.location.reload()
				}
			} else {
				// show error message
				errorNotification({
					message: 'Organization creation failed'
				})
			}
		} catch (error) {
			console.error('error', error)
			errorNotification({
				message: 'Organization creation failed'
			})
		}
	}

	async function updateOrganizationWhatsAppBusinessAccountDetails(
		data: z.infer<typeof WhatsappBusinessAccountDetailsFormSchema>
	) {
		try {
			if (!currentOrganization) return

			const response = await updateWhatsappBusinessAccountDetailsMutation.mutateAsync({
				data: {
					businessAccountId: data.whatsappBusinessAccountId,
					accessToken: data.apiToken,
					webhookSecret: data.webhookSecret
				}
			})

			if (response.accessToken) {
				writeProperty({
					currentOrganization: {
						...currentOrganization,
						whatsappBusinessAccountDetails: {
							businessAccountId: data.whatsappBusinessAccountId,
							accessToken: data.apiToken,
							webhookSecret: data.webhookSecret
						}
					}
				})

				router.push('/onboarding/invite-team-members')
			} else {
				errorNotification({
					message: 'Error updating WhatsApp Business Account ID'
				})
			}
		} catch (error) {
			console.error(error)
			errorNotification({
				message: 'Error updating WhatsApp Business Account ID'
			})
		}
	}

	switch (step.slug) {
		case 'create-organization': {
			return (
				<div className="flex w-full items-center justify-end space-x-2">
					<Form {...newOrganizationForm}>
						<form
							onSubmit={newOrganizationForm.handleSubmit(handleCreateOrganization)}
							className="w-full space-y-8"
						>
							<div className="flex flex-col gap-8">
								<FormField
									control={newOrganizationForm.control}
									name="name"
									render={({ field }) => (
										<FormItem className="w-full">
											<FormLabel>Name</FormLabel>
											<FormControl>
												<Input
													placeholder="name"
													{...field}
													autoComplete="off"
													className="w-full"
												/>
											</FormControl>
											<FormMessage />
										</FormItem>
									)}
								/>

								<FormField
									control={newOrganizationForm.control}
									name="description"
									render={({ field }) => (
										<FormItem>
											<FormLabel>Description</FormLabel>
											<FormControl>
												<Textarea
													placeholder="Description (optional)"
													{...field}
													autoComplete="off"
												/>
											</FormControl>
											<FormMessage />
										</FormItem>
									)}
								/>
							</div>
							<Button className="ml-auto mr-0 w-full" type="submit">
								Create Organization
							</Button>
						</form>
					</Form>
				</div>
			)
		}

		case 'invite-team-members': {
			return (
				<div className="flex w-full items-center justify-end space-x-2 pt-6">
					<Form {...newMemberInviteForm}>
						<form
							onSubmit={newMemberInviteForm.handleSubmit(inviteTeamMembers)}
							className="w-full space-y-8"
						>
							<div className="flex flex-col gap-8">
								<FormField
									control={newMemberInviteForm.control}
									name="email"
									render={({ field }) => (
										<FormItem>
											<FormLabel>Email</FormLabel>
											<FormControl>
												<Input
													disabled={isBusy}
													placeholder="Email"
													{...field}
													autoComplete="off"
												/>
											</FormControl>
											<FormMessage />
										</FormItem>
									)}
								/>
							</div>
							<Button disabled={isBusy} className="ml-auto mr-0 w-full" type="submit">
								Invite Now
							</Button>

							<p
								className="cursor-pointer text-center text-sm text-gray-400 hover:underline"
								onClick={() => {
									router.push('/dashboard')
								}}
							>
								Skip to dashboard
							</p>
						</form>
					</Form>
				</div>
			)
		}

		case 'whatsapp-business-account-details': {
			return (
				<div className="flex w-full max-w-4xl items-center justify-end space-x-2">
					<Form {...whatsappBusinessAccountForm}>
						<form
							onSubmit={whatsappBusinessAccountForm.handleSubmit(
								updateOrganizationWhatsAppBusinessAccountDetails
							)}
							className="flex w-full flex-col gap-4"
						>
							<FormField
								control={whatsappBusinessAccountForm.control}
								name="whatsappBusinessAccountId"
								render={({ field }) => (
									<FormItem className="w-full">
										<FormLabel>WhatsApp Business Account ID</FormLabel>
										<FormControl>
											<div className="flex flex-row gap-2">
												<Input
													disabled={isBusy}
													placeholder="whatsapp business account id"
													{...field}
													autoComplete="off"
													type={
														whatsAppBusinessAccountDetailsVisibility.whatsappBusinessAccountId
															? 'text'
															: 'password'
													}
												/>
												<span
													className="rounded-md border p-1 px-2"
													onClick={() => {
														setWhatsAppBusinessAccountDetailsVisibility(
															data => ({
																...data,
																whatsappBusinessAccountId:
																	!data.whatsappBusinessAccountId
															})
														)
													}}
												>
													<EyeIcon className="size-5" />
												</span>
											</div>
										</FormControl>
										<FormMessage />
									</FormItem>
								)}
							/>
							<FormField
								control={whatsappBusinessAccountForm.control}
								name="apiToken"
								render={({ field }) => (
									<FormItem className="w-full">
										<FormLabel>API key</FormLabel>
										<FormControl>
											<div className="flex flex-row gap-2">
												<Input
													disabled={isBusy}
													placeholder="whatsapp business account api token"
													{...field}
													autoComplete="off"
													type={
														whatsAppBusinessAccountDetailsVisibility.apiToken
															? 'text'
															: 'password'
													}
												/>
												<span
													className="rounded-md border p-1 px-2"
													onClick={() => {
														setWhatsAppBusinessAccountDetailsVisibility(
															data => ({
																...data,
																apiToken: !data.apiToken
															})
														)
													}}
												>
													<EyeIcon className="size-5" />
												</span>
											</div>
										</FormControl>
										<FormMessage />
									</FormItem>
								)}
							/>
							<FormField
								control={whatsappBusinessAccountForm.control}
								name="webhookSecret"
								render={({ field }) => (
									<FormItem>
										<FormLabel>Webhook Secret</FormLabel>
										<FormControl>
											<div className="flex flex-row gap-2">
												<Input
													disabled={isBusy}
													placeholder="whatsapp business account webhook secret"
													{...field}
													autoComplete="off"
													type={
														whatsAppBusinessAccountDetailsVisibility.webhookSecret
															? 'text'
															: 'password'
													}
												/>
												<span
													className="rounded-md border p-1 px-2"
													onClick={() => {
														setWhatsAppBusinessAccountDetailsVisibility(
															data => ({
																...data,
																webhookSecret: !data.webhookSecret
															})
														)
													}}
												>
													<EyeIcon className="size-5" />
												</span>
											</div>
										</FormControl>
										<FormMessage />
									</FormItem>
								)}
							/>
							<Button
								type="submit"
								className="ml-auto w-full"
								disabled={isBusy || !whatsappBusinessAccountForm.formState.isDirty}
							>
								Update
							</Button>
						</form>
					</Form>
				</div>
			)
		}
	}
}

export default OnboardingStepClientPage
