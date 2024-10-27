'use client'

import { type ChangeEvent, useRef } from 'react'
import { Button } from '~/components/ui/button'

export const FileUploaderComponent: React.FC<{
	descriptionString: string
	onFileUpload:
		| ((e: ChangeEvent<HTMLInputElement>) => Promise<void>)
		| ((e: ChangeEvent<HTMLInputElement>) => void)
}> = ({ descriptionString, onFileUpload }) => {
	const fileInputRef = useRef<HTMLInputElement>(null)

	return (
		<div className="flex flex-col items-center justify-center gap-4 rounded-lg border border-dashed p-5 text-center">
			{descriptionString ? (
				<label className="text-wrap text-sm font-semibold">{descriptionString}</label>
			) : null}

			<input
				ref={fileInputRef}
				hidden
				type="file"
				id="file_uploader_input"
				onChange={e => {
					onFileUpload(e)
				}}
			/>

			<Button
				type="button"
				size={'sm'}
				onClick={() => {
					console.log('trigger input file event')
					if (fileInputRef.current) {
						fileInputRef.current.click()
					} else {
						console.warn('fileInputRef.current is null!!!!')
					}
				}}
			>
				Upload
			</Button>
		</div>
	)
}
