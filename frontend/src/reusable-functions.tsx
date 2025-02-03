import { v4 } from 'uuid'
import { toast } from 'sonner'
import { CheckCircledIcon, InfoCircledIcon } from '@radix-ui/react-icons'
import { createRoot } from 'react-dom/client'
import { AlertModal } from './components/modal/alert-modal'
import { type MessageTemplateSchema } from 'root/.generated'
import { Icons } from './components/icons'

export function generateUniqueId() {
	return v4()
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
			<Icons.xCircle className="h-5 w-5" color="#ef4444" />
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
			<Icons.warning className="h-5 w-5" color="#fcb603" />
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

/**
 * Calculates the number of parameters required for each component in the template.
 * @param template - The template object.
 * @returns An object with component types as keys and parameter counts as values.
 */
export function getParametersPerComponent(
	template?: MessageTemplateSchema
): Record<string, number> {
	const parameterCounts: Record<string, number> = {}

	if (!template || !template.components) {
		return parameterCounts
	}

	template.components.forEach(component => {
		if (!component.type) {
			return
		}

		let parameterCount = 0

		// Check the example field of the main component
		if (component.example) {
			switch (component.type) {
				case 'BODY': {
					if (component.example.body_text?.length) {
						// it is an array of array
						component.example.body_text.forEach(bodyText => {
							parameterCount += bodyText.length
						})
					}
					break
				}

				case 'HEADER': {
					if (component.example.header_text) {
						parameterCount += component.example.header_text.length
					}

					if (component.format !== 'TEXT') {
						parameterCount += 1
					}

					break
				}
			}
		}

		// Check the example field of any buttons
		// ! TODO: enable this after fixing the wapi.go object structure for template buttons
		if (component.buttons) {
			component.buttons.forEach(button => {
				if (button.example) {
					parameterCount += button.example.length
				}
			})
		}

		const keyToUse =
			component.type === 'BODY' ? 'body' : component.type === 'BUTTONS' ? 'buttons' : 'header'

		// Add the count for this component
		parameterCounts[keyToUse] = parameterCount
	})

	return parameterCounts
}
