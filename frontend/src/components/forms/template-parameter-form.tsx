import { type UseFormReturn } from 'react-hook-form'
import { Button } from '~/components/ui/button'
import { type TemplateComponentSchema } from '~/schema'
import { type z } from 'zod'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '../ui/form'
import { getParametersPerComponent } from '~/reusable-functions'
import { type MessageTemplateComponentType, type MessageTemplateSchema } from 'root/.generated'
import { Input } from '../ui/input'
import { Separator } from '../ui/separator'

type Props = {
	isBusy: boolean
	setIsTemplateComponentsInputModalOpen: (isOpen: boolean) => void
	templateMessageComponentParameterForm: UseFormReturn<
		{
			body: (string | null | undefined)[]
			header: (string | null | undefined)[]
			buttons: (string | null | undefined)[]
		},
		any,
		undefined
	>
	handleTemplateComponentParameterSubmit: (
		data: z.infer<typeof TemplateComponentSchema>
	) => Promise<void>
	template?: MessageTemplateSchema
}

const TemplateParameterForm: React.FC<Props> = ({
	isBusy,
	setIsTemplateComponentsInputModalOpen,
	templateMessageComponentParameterForm,
	handleTemplateComponentParameterSubmit,
	template
}) => {
	console.log({
		templateMessageComponentParameterFormValues: JSON.stringify(
			templateMessageComponentParameterForm.watch()
		)
	})

	console.log({
		isValidating: templateMessageComponentParameterForm.formState.isValidating,
		isSubmitting: templateMessageComponentParameterForm.formState.isSubmitting
	})

	return (
		<Form {...templateMessageComponentParameterForm}>
			<form
				onSubmit={templateMessageComponentParameterForm.handleSubmit(
					handleTemplateComponentParameterSubmit
				)}
				className="flex-1 space-y-8 px-2 "
			>
				<div className="flex max-h-[32rem] flex-col gap-8 overflow-scroll pb-44">
					{Object.entries(getParametersPerComponent(template)).map(
						([key, value], index) => {
							if (!value) {
								return null
							}

							const componentType = key as Lowercase<MessageTemplateComponentType>

							const component = template?.components?.find(component => {
								return (
									component.type &&
									component.type === (componentType.toUpperCase() as any)
								)
							})

							return (
								<div key={`${key}_parameters`} className="flex flex-col gap-3">
									<span className="font-bold">{key}</span>

									{key === 'buttons' ? (
										<div
											key={'buttons-parameters'}
											className="flex w-full flex-col gap-3"
										>
											{component?.buttons?.map((button, buttonIndex) => {
												if (
													button.example?.length ||
													button.type === 'QUICK_REPLY'
												) {
													let exampleValue = button.example?.[0]
													let placeHolderSuffix = 'parameter'

													if (button.type === 'QUICK_REPLY') {
														exampleValue = 'PAYLOAD'
														placeHolderSuffix = 'payload'
													}

													return (
														<FormField
															key={`${componentType}-${buttonIndex}`}
															control={
																templateMessageComponentParameterForm.control
															}
															name={'buttons'}
															render={({ field }) => (
																<FormItem>
																	<FormLabel>
																		<span className="flex flex-row">
																			{button.type}{' '}
																			{placeHolderSuffix}
																			&nbsp;
																			<pre className="italic text-red-500">
																				Example:{' '}
																				{exampleValue}
																			</pre>
																		</span>
																	</FormLabel>
																	<FormControl>
																		<Input
																			disabled={isBusy}
																			placeholder={`${button.type} - ${placeHolderSuffix}`}
																			{...field}
																			autoComplete="off"
																			value={
																				(field?.value &&
																					field.value[
																						buttonIndex
																					]) ||
																				''
																			}
																			onChange={e => {
																				// existing params
																				const existingParamValue =
																					templateMessageComponentParameterForm.getValues(
																						'buttons'
																					)

																				if (
																					existingParamValue
																				) {
																					existingParamValue[
																						buttonIndex
																					] =
																						e.target.value

																					templateMessageComponentParameterForm.setValue(
																						'buttons',
																						existingParamValue
																					)
																				} else {
																					// create a new object
																					const paramArray =
																						[]

																					paramArray[
																						buttonIndex
																					] =
																						e.target.value

																					console.log({
																						paramArray
																					})

																					templateMessageComponentParameterForm.setValue(
																						'buttons',
																						paramArray
																					)
																				}
																			}}
																		/>
																	</FormControl>
																	<FormMessage />
																</FormItem>
															)}
														/>
													)
												} else {
													return null
												}
											})}
										</div>
									) : (
										<>
											{Array(value)
												.fill(0)
												.map((_, index) => {
													let exampleValue: string | null = null

													if (component) {
														switch (componentType) {
															case 'header':
																exampleValue =
																	component.example
																		?.header_text?.[index] ||
																	null
																break
															case 'body':
																exampleValue =
																	component.example
																		?.body_text?.[0][index] ||
																	null
																break
															default:
																exampleValue = null
																break
														}
													}

													return (
														<FormField
															key={`${componentType}-${index}`}
															control={
																templateMessageComponentParameterForm.control
															}
															name={
																componentType as 'header' | 'body'
															}
															render={({ field }) => (
																<FormItem>
																	<FormLabel>
																		<span className="flex flex-row">
																			{`{{${index + 1}}}`}
																			&nbsp;
																			<pre className="italic text-red-500">
																				Example:{' '}
																				{exampleValue}
																			</pre>
																		</span>
																	</FormLabel>
																	<FormControl>
																		<Input
																			disabled={isBusy}
																			placeholder={`${componentType} - {{${index + 1}}} - ${exampleValue}`}
																			{...field}
																			autoComplete="off"
																			value={
																				(field?.value &&
																					field.value[
																						index
																					]) ||
																				''
																			}
																			onChange={e => {
																				// existing params

																				const existingParamValue =
																					templateMessageComponentParameterForm.getValues(
																						componentType as
																							| 'header'
																							| 'body'
																					)

																				if (
																					existingParamValue
																				) {
																					existingParamValue[
																						index
																					] =
																						e.target.value

																					templateMessageComponentParameterForm.setValue(
																						componentType as
																							| 'header'
																							| 'body',
																						existingParamValue
																					)
																				} else {
																					// create a new object
																					const paramArray =
																						[]

																					paramArray[
																						index
																					] =
																						e.target.value

																					templateMessageComponentParameterForm.setValue(
																						componentType as
																							| 'header'
																							| 'body'
																							| 'buttons',
																						paramArray
																					)
																				}
																			}}
																		/>
																	</FormControl>
																	<FormMessage />
																</FormItem>
															)}
														/>
													)
												})}
										</>
									)}

									{index < 2 && <Separator className="mt-6" />}
								</div>
							)
						}
					)}
				</div>

				<div className="sticky bottom-36 flex w-full flex-col gap-3 bg-background py-10">
					<pre className="text-xs text-red-500">NOTE: Scroll for more inputs</pre>
					<div className="flex w-full flex-row gap-3">
						<Button disabled={isBusy} className="ml-auto mr-0 w-full" type="submit">
							Save
						</Button>
						<Button
							disabled={isBusy}
							variant={'outline'}
							className="ml-auto mr-0 w-full"
							type="button"
							onClick={() => {
								setIsTemplateComponentsInputModalOpen(false)
							}}
						>
							Cancel
						</Button>
					</div>
				</div>
			</form>
		</Form>
	)
}

export default TemplateParameterForm
