'use client'

import { type ChangeEvent, type HTMLAttributes, useRef } from 'react'
import { Button } from '~/components/ui/button'

export const FileUploaderComponent: React.FC<
	{
		descriptionString: string
		onFileUpload: (e: ChangeEvent<HTMLInputElement>) => void
	} & HTMLAttributes<HTMLInputElement>
> = ({ descriptionString, onFileUpload, ...props }) => {
	const fileInputRef = useRef<HTMLInputElement | null>(null)

	return (
		<div className="flex flex-col items-center justify-center gap-4 rounded-lg border border-dashed p-5 text-center">
			{descriptionString ? (
				<label className="text-wrap text-sm font-semibold">{descriptionString}</label>
			) : null}

			<input
				hidden
				type="file"
				id="file_uploader_input"
				accept=".csv"
				{...props}
				ref={fileInputRef}
				onChange={e => {
					onFileUpload(e)
				}}
			/>

			{fileInputRef.current?.files?.length ? (
				<>
					<div className="flex flex-row items-center gap-2">
						<label className="text-sm font-semibold"></label>
						<p className="text-sm">{fileInputRef.current.files[0].name}</p>
						{/* <span
							onClick={() => {
								if (fileInputRef.current) {
									fileInputRef.current.files = null
								}
							}}
							className="cursor-pointer"
						>
							<XCircle className="size-5 text-secondary-foreground" />
						</span> */}
					</div>
				</>
			) : null}

			<Button
				type="button"
				size={'sm'}
				onClick={() => {
					if (fileInputRef.current) {
						fileInputRef.current.click()
					} else {
						console.warn('fileInputRef.current is null!!!!')
					}
				}}
			>
				{fileInputRef.current?.files?.length ? 'Change File' : 'Upload File'}
			</Button>
		</div>
	)
}
