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
import { useForm } from 'react-hook-form'
import { type z } from 'zod'
import { errorNotification, successNotification } from '~/reusable-functions'
import { BulkImportContactsFormSchema } from '~/schema'
import { useCallback, useState } from 'react'
import { useGetContactLists } from 'root/.generated'
import { FileUploaderComponent } from '~/components/file-uploader'
import { MultiSelect } from '~/components/multi-select'
import { AUTH_TOKEN_LS, getBackendUrl } from '~/constants'
import DocumentationPitch from '~/components/forms/documentation-pitch'
import BreadCrumb from '~/components/breadcrumb'
import { Heading } from '~/components/ui/heading'
import { Separator } from '~/components/ui/separator'
import Progress from '~/components/progress'
import { useScrollToBottom } from '~/hooks/use-scroll-to-bottom'

const BulkImportContactPage = () => {
	const bulkImportForm = useForm<z.infer<typeof BulkImportContactsFormSchema>>({
		resolver: zodResolver(BulkImportContactsFormSchema)
	})

	const [messagesContainerRef, messagesEndRef] = useScrollToBottom<HTMLDivElement>()
	const [abortController, setAbortController] = useState<AbortController | null>(null)
	const [importState, setImportState] = useState<'importing' | 'done' | 'idle' | 'error'>('idle')

	const [isBusy, setIsBusy] = useState(false)
	const [file, setFile] = useState<File | null>(null)
	const [progress, setProgress] = useState({ current: 0, total: 0 })
	const [logs, setLogs] = useState<string[]>([])

	const listsResponse = useGetContactLists({
		order: 'asc',
		page: 1,
		per_page: 50
	})

	const handleDataStream = useCallback(async (reader: ReadableStreamDefaultReader) => {
		const decoder = new TextDecoder()
		let buffer = ''

		let done = false
		while (!done) {
			const { done: readerDone, value } = await reader.read()
			done = readerDone
			if (done) break

			buffer += decoder.decode(value, { stream: true })
			const chunks = buffer.split('\n')

			// Process each chunk
			for (let i = 0; i < chunks.length - 1; i++) {
				const chunk = chunks[i]
				const data = JSON.parse(chunk)

				switch (data.type) {
					case 'importing':
						setImportState(() => 'importing')
						setLogs(l => [...l, data.message])
						break
					case 'progress':
						setProgress(p => ({
							current: data.current,
							total: data.total || p.total
						}))
						setLogs(l => [...l, data.message])
						break
					case 'error':
						setLogs(l => [...l, `Error: ${data.message}`])
						setImportState(() => 'error')
						setIsBusy(false)
						break
					case 'complete':
						setLogs(l => [...l, `Complete: ${data.message}`])
						setIsBusy(false)
						setImportState(() => 'done')
						successNotification({ message: data.message })
						break
				}
			}

			// Keep the last incomplete chunk in the buffer
			buffer = chunks[chunks.length - 1]
		}
	}, [])

	async function onBulkContactImportFormSubmit(
		data: z.infer<typeof BulkImportContactsFormSchema>
	) {
		try {
			setIsBusy(true)
			setLogs([])
			setProgress({ current: 0, total: 0 })

			if (!file) {
				errorNotification({
					message: 'Please upload a file'
				})
				return
			}

			const controller = new AbortController()
			setAbortController(controller)

			const formData = new FormData()
			formData.append('file', file)
			formData.append('delimiter', data.delimiter)
			formData.append('listIds', JSON.stringify(data.listIds)) // If `listIds` is an array, stringify it

			const response = await fetch(`${getBackendUrl()}/contacts/bulkImport`, {
				body: formData,
				method: 'POST',
				headers: {
					Accept: 'application/json',
					'x-access-token': localStorage.getItem(AUTH_TOKEN_LS) || ''
				},
				cache: 'no-cache',
				credentials: 'include',
				signal: controller.signal
			})

			if (!response.body) throw new Error('No response body')
			const reader = response.body.getReader()

			await handleDataStream(reader)
		} catch (error) {
			console.error(error)
			errorNotification({
				message: 'An error occurred'
			})
		} finally {
			setIsBusy(false)
		}
	}

	const breadcrumbItems = [
		{ title: 'Contacts', link: '/contacts' },
		{ title: 'Bulk Import', link: '/contacts/bulk-import' }
	]

	return (
		<>
			<div className="flex-1 space-y-4  p-4 pt-6 md:p-8">
				<BreadCrumb items={breadcrumbItems} />
				<div className="flex items-start justify-between">
					<Heading title={`Bulk import`} description="" />
				</div>
				<Separator />

				<div className="flex flex-row gap-10">
					{importState === 'idle' ? (
						<div className="flex w-full items-center justify-end space-x-2">
							<Form {...bulkImportForm}>
								<form
									onSubmit={bulkImportForm.handleSubmit(
										onBulkContactImportFormSubmit
									)}
									className="w-full space-y-8"
								>
									<div className="flex flex-col gap-8">
										<FormField
											control={bulkImportForm.control}
											name="file"
											render={({ field }) => (
												<FormItem>
													<FormLabel>Upload CSV File</FormLabel>
													<FileUploaderComponent
														descriptionString="CSV File"
														{...field}
														onFileUpload={e => {
															const file = e.target.files?.[0]

															if (!file) return
															setFile(() => file)
														}}
													/>
												</FormItem>
											)}
										/>

										<FormField
											control={bulkImportForm.control}
											name="delimiter"
											render={({ field }) => (
												<FormItem>
													<FormLabel>Delimiter</FormLabel>
													<FormControl>
														<Input
															disabled={isBusy}
															placeholder="Column delimiter (e.g. ,)"
															{...field}
															autoComplete="off"
														/>
													</FormControl>
													<FormMessage />
												</FormItem>
											)}
										/>
										<FormField
											control={bulkImportForm.control}
											name="listIds"
											render={({}) => (
												<FormItem className="tablet:w-3/4 tablet:gap-2 desktop:w-1/2 flex flex-col gap-1 ">
													<FormLabel>Select the lists</FormLabel>
													<MultiSelect
														options={
															listsResponse?.data?.lists.map(
																list => ({
																	label: list.name,
																	value: list.uniqueId
																})
															) || []
														}
														onValueChange={e => {
															bulkImportForm.setValue('listIds', e, {
																shouldValidate: true
															})
														}}
														defaultValue={bulkImportForm.watch(
															'listIds'
														)}
														placeholder="Select lists"
														variant="default"
													/>
													<FormMessage />
												</FormItem>
											)}
										/>
									</div>
									<Button
										disabled={isBusy}
										className="ml-auto mr-0 w-full"
										type="submit"
									>
										Import
									</Button>
								</form>
							</Form>
						</div>
					) : null}

					{importState === 'importing' || importState === 'error' ? (
						<div className="flex w-full  flex-1 flex-col gap-10 space-x-2">
							{/* Progress bar */}
							<div className="mb-4">
								<Progress value={(progress.current / progress.total) * 100} />
								<div className="mt-1 text-sm text-muted-foreground">
									{progress.current} of {progress.total} processed
								</div>
							</div>

							{/* Logs container */}
							<div className="mb-4 h-64 overflow-y-auto rounded border p-2">
								{logs.map((log, i) => (
									<div
										key={i}
										className="font-mono text-sm"
										ref={messagesContainerRef}
									>
										{log}
										<div
											ref={messagesEndRef}
											className="min-h-[24px] min-w-[24px] shrink-0"
										/>
									</div>
								))}
							</div>

							{/* Cancel button */}
							<Button
								variant="destructive"
								onClick={() => {
									if (abortController) {
										abortController.abort()
									}
									setIsBusy(false)
									setImportState('idle')
								}}
							>
								Cancel Import
							</Button>
						</div>
					) : null}

					<DocumentationPitch type="contact" />
				</div>
			</div>
		</>
	)
}

export default BulkImportContactPage
