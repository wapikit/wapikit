'use client'
import { useEffect, useState } from 'react'
import { Button } from '~/components/ui/button'
import { Modal } from '~/components/ui/modal'

interface AlertModalProps {
	isOpen: boolean
	onClose: () => void
	onConfirm: (confirmation: boolean) => void
	loading: boolean
	title: string
	description: string
}

export const AlertModal: React.FC<AlertModalProps> = ({
	isOpen,
	onClose,
	onConfirm,
	loading,
	title,
	description
}) => {
	const [isMounted, setIsMounted] = useState(false)

	useEffect(() => {
		setIsMounted(true)
	}, [])

	if (!isMounted) {
		return null
	}

	return (
		<Modal
			title={title}
			description={description}
			isOpen={isOpen}
			onClose={onClose}
			isDismissible={false}
		>
			<div className="flex w-full items-center justify-end space-x-2 pt-6">
				<Button
					disabled={loading}
					variant="outline"
					onClick={() => {
						onConfirm(false)
						onClose()
					}}
				>
					Cancel
				</Button>
				<Button
					disabled={loading}
					variant="destructive"
					onClick={() => {
						onConfirm(true)
						onClose()
					}}
				>
					Continue
				</Button>
			</div>
		</Modal>
	)
}
