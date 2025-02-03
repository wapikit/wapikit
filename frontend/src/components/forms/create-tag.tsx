'use client'

import { useForm } from 'react-hook-form'
import { Modal } from '../ui/modal'
import { type z } from 'zod'
import { CreateTagFormSchema } from '~/schema'
import { zodResolver } from '@hookform/resolvers/zod'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '../ui/form'
import { Input } from '../ui/input'
import { Button } from '../ui/button'
import { useCreateOrganizationTag } from 'root/.generated'
import { useState } from 'react'
import { errorNotification } from '~/reusable-functions'
import { useLayoutStore } from '~/store/layout.store'

const CreateTagModal = () => {
	const { isCreateTagModalOpen, writeProperty, tags: existingTags } = useLayoutStore()

	const [isBusy, setIsBusy] = useState(false)

	const createTagForm = useForm<z.infer<typeof CreateTagFormSchema>>({
		resolver: zodResolver(CreateTagFormSchema),
		defaultValues: {
			label: ''
		}
	})

	const createTagMutation = useCreateOrganizationTag()

	async function createTag(data: z.infer<typeof CreateTagFormSchema>) {
		try {
			if (isBusy) return

			setIsBusy(true)

			const res = await createTagMutation.mutateAsync({
				data: {
					label: data.label
				}
			})

			if (res) {
				writeProperty({
					isCreateTagModalOpen: false,
					tags: [...existingTags, res.tag]
				})
			}
		} catch (error) {
			console.error(error)
			errorNotification({
				message: 'Failed to create tag'
			})
		} finally {
			setIsBusy(false)
		}
	}

	return (
		<Modal
			title="Create a new tag"
			description=""
			isOpen={isCreateTagModalOpen}
			onClose={() => {
				writeProperty({
					isCreateTagModalOpen: false
				})
			}}
		>
			<div className="flex w-full items-center justify-end space-x-2 pt-6">
				<Form {...createTagForm}>
					<form
						onSubmit={createTagForm.handleSubmit(createTag)}
						className="w-full space-y-8"
					>
						<div className="flex flex-col gap-8">
							<FormField
								control={createTagForm.control}
								name="label"
								render={({ field }) => (
									<FormItem>
										<FormLabel>Tag Label</FormLabel>
										<FormControl>
											<Input
												disabled={isBusy}
												placeholder="Label"
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
							Create Tag
						</Button>
					</form>
				</Form>
			</div>
		</Modal>
	)
}

export default CreateTagModal
