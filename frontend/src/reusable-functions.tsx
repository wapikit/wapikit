import { ulid } from 'ulid'
import { toast } from 'sonner'
import { MessageSquareWarning, XCircleIcon } from 'lucide-react'
import { CheckCircledIcon, InfoCircledIcon } from '@radix-ui/react-icons'
import { createRoot } from 'react-dom/client'
import { AlertModal } from './components/modal/alert-modal'

export function generateUniqueId() {
	return ulid()
}

export function getWebsocketUrl(token: string) {
	return process.env.NODE_ENV === 'development' ? `ws://127.0.0.1:8081/ws?token=${token}` : ``
}

export function infoNotification(params: { message: string; darkMode?: true; duration?: string }) {
	return toast.error(
		<div className="flex flex-row items-center justify-start gap-2">
			<InfoCircledIcon className="h-5 w-5" color="#3b82f6" />
			<span>{params.message}</span>
		</div>
	)
}

export function errorNotification(params: { message: string; darkMode?: true; duration?: string }) {
	return toast.error(
		<div className="flex flex-row items-center justify-start gap-2">
			<XCircleIcon className="h-5 w-5" color="#ef4444" />
			<span>{params.message}</span>
		</div>
	)
}

export function successNotification(params: {
	message: string
	darkMode?: true
	duration?: string
}) {
	return toast.success(
		<div className="flex flex-row items-center justify-start gap-2">
			<CheckCircledIcon className="h-5 w-5" color="#22c55e" />
			<span>{params.message}</span>
		</div>
	)
}

export function warnNotification(params: { message: string; darkMode?: true; duration?: string }) {
	return toast.success(
		<div className="flex flex-row items-center justify-start gap-2">
			<MessageSquareWarning className="h-5 w-5" color="#fcb603" />
			<span>{params.message}</span>
		</div>
	)
}

export function materialConfirm(params: { title: string; description: string }): Promise<boolean> {
	return new Promise(resolve => {
		const container = document.createElement('div')
		document.body.appendChild(container)

		const root = createRoot(container)

		const closeDialog = () => {
			root.unmount()
		}

		const handleConfirm = (confirmation: boolean) => {
			closeDialog()
			resolve(confirmation)
		}

		root.render(
			<AlertModal
				loading={false}
				isOpen={true}
				onClose={closeDialog}
				onConfirm={handleConfirm}
				title={params.title}
				description={params.description}
			/>
		)
	})
}

export function parseMessageContentForHyperLink(message: string) {
	const urlRegex = /(https?:\/\/[^\s]+)/g
	return message.replace(urlRegex, '<a href="$1" target="_blank">$1</a>')
}
