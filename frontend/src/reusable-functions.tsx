import { ulid } from 'ulid'
import { toast } from 'sonner'
import { MessageSquareWarning, XCircleIcon } from 'lucide-react'
import { CheckCircledIcon, InfoCircledIcon } from '@radix-ui/react-icons'

export function generateUniqueId() {
	return ulid()
}

export function getWebsocketUrl(token: string) {
	return process.env.NODE_ENV === 'development' ? `ws://localhost:3001?token=${token}` : ``
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

export function materialConfirm() {}
