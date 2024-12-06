import { type MessageTemplateSchema } from 'root/.generated'
import { Icons } from '../icons'
import { type z } from 'zod'
import { type TemplateComponentSchema } from '~/schema'
import { clsx } from 'clsx'
import dayjs from 'dayjs'
import { Separator } from '../ui/separator'
import MessageButtonRenderer from './button-render'

const TemplateMessageRenderer: React.FC<{
	templateMessage?: MessageTemplateSchema
	parameterValues: z.infer<typeof TemplateComponentSchema>
}> = ({ templateMessage, parameterValues }) => {
	if (!templateMessage) {
		return null
	}

	const header = templateMessage.components?.find(component => component.type === 'HEADER')
	const body = templateMessage.components?.find(component => component.type === 'BODY')
	const footer = templateMessage.components?.find(component => component.type === 'FOOTER')
	const buttons = templateMessage.components?.find(
		component => component.type === 'BUTTONS'
	)?.buttons

	let headerText = header?.text || ''
	let bodyText = body?.text || ''

	if (headerText.includes('{{')) {
		headerText = headerText.replace(/{{(.*?)}}/g, (_, match: string) => {
			const index = Number(match) - 1
			const value = parameterValues.header?.[index]
			return value || `{{${match}}}`
		})
	}

	if (bodyText.includes('{{')) {
		bodyText = bodyText.replace(/{{(.*?)}}/g, (_, match: string) => {
			const index = Number(match) - 1
			const value = parameterValues.body?.[index]
			return value || `{{${match}}}`
		})
	}

	const MenuIcon = Icons.menu

	return (
		<div
			className={clsx(
				'mr-auto flex   max-w-96 flex-col gap-2 rounded-md bg-white p-1 px-3 text-foreground'
			)}
		>
			{/* header */}
			{header ? <p className="font-bold">{headerText}</p> : null}

			{/* body */}
			{body ? <p className="text-sm">{bodyText}</p> : null}

			{/* footer */}
			<div className="flex flex-row items-start justify-between gap-1">
				{footer ? <p className="flex-1 text-xs opacity-75">{footer.text}</p> : null}
				<span className={clsx('ml-auto text-[10px]')}>{dayjs().format('hh:mm A')}</span>
			</div>

			{/*  buttons */}
			{buttons?.length ? (
				<div>
					<Separator className="w-full" />
					{buttons.map((button, index) => {
						if (index > 1) {
							return null
						}

						return (
							<>
								<MessageButtonRenderer
									key={`${button.type}-${button.text}`}
									messageButton={button}
								/>

								{index === buttons.length - 1 ? null : (
									<Separator key={`${index}-separator`} />
								)}
							</>
						)
					})}

					{buttons.length > 2 ? (
						<div className="flex cursor-pointer items-center  justify-center gap-2 py-2 text-center text-blue-500">
							<MenuIcon className="size-5" />
							See All Options
						</div>
					) : null}
				</div>
			) : null}
		</div>
	)
}

export default TemplateMessageRenderer
