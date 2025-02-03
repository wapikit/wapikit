import {
	ArrowRight,
	CircuitBoardIcon,
	Command,
	CreditCard,
	HelpCircle,
	Loader2,
	LogIn,
	type LucideIcon,
	type LucideProps,
	MoreVertical,
	Reply,
	Forward,
	Phone,
	CheckCheck,
	Pointer,
	Brain
} from 'lucide-react'
import {
	DashboardIcon,
	CalendarIcon,
	CodeIcon,
	BellIcon,
	ExternalLinkIcon,
	DotsVerticalIcon,
	InfoCircledIcon,
	ClipboardCopyIcon,
	Cross2Icon,
	ChevronLeftIcon,
	ChevronRightIcon,
	TrashIcon,
	ImageIcon,
	GearIcon,
	FileIcon,
	FileTextIcon,
	PlusIcon,
	MoonIcon,
	SunIcon,
	Link2Icon,
	CrossCircledIcon,
	PauseIcon,
	PlayIcon,
	ExclamationTriangleIcon,
	RocketIcon,
	ChatBubbleIcon,
	PersonIcon,
	TwitterLogoIcon,
	RowsIcon,
	MinusCircledIcon,
	CheckIcon,
	Pencil2Icon,
	DownloadIcon
} from '@radix-ui/react-icons'
import { Avatar } from '@radix-ui/react-avatar'

export type Icon = LucideIcon

export const Icons = {
	download: DownloadIcon,
	brain: Brain,
	bell: BellIcon,
	pointer: Pointer,
	calendar: CalendarIcon,
	jsonBrackets: CodeIcon,
	doubleCheck: CheckCheck,
	externalLink: ExternalLinkIcon,
	phone: Phone,
	menu: DotsVerticalIcon,
	info: InfoCircledIcon,
	forward: Forward,
	reply: Reply,
	clipboard: ClipboardCopyIcon,
	dashboard: DashboardIcon,
	logo: Command,
	login: LogIn,
	close: Cross2Icon,
	profile: Avatar,
	spinner: Loader2,
	kanban: CircuitBoardIcon,
	chevronLeft: ChevronLeftIcon,
	chevronRight: ChevronRightIcon,
	trash: TrashIcon,
	employee: Avatar,
	post: FileTextIcon,
	page: FileIcon,
	sparkles: ({ size = 16 }: { size?: number }) => (
		<svg
			height={size}
			strokeLinejoin="round"
			viewBox="0 0 16 16"
			width={size}
			style={{ color: 'currentcolor' }}
		>
			<path
				d="M2.5 0.5V0H3.5V0.5C3.5 1.60457 4.39543 2.5 5.5 2.5H6V3V3.5H5.5C4.39543 3.5 3.5 4.39543 3.5 5.5V6H3H2.5V5.5C2.5 4.39543 1.60457 3.5 0.5 3.5H0V3V2.5H0.5C1.60457 2.5 2.5 1.60457 2.5 0.5Z"
				fill="currentColor"
			/>
			<path
				d="M14.5 4.5V5H13.5V4.5C13.5 3.94772 13.0523 3.5 12.5 3.5H12V3V2.5H12.5C13.0523 2.5 13.5 2.05228 13.5 1.5V1H14H14.5V1.5C14.5 2.05228 14.9477 2.5 15.5 2.5H16V3V3.5H15.5C14.9477 3.5 14.5 3.94772 14.5 4.5Z"
				fill="currentColor"
			/>
			<path
				d="M8.40706 4.92939L8.5 4H9.5L9.59294 4.92939C9.82973 7.29734 11.7027 9.17027 14.0706 9.40706L15 9.5V10.5L14.0706 10.5929C11.7027 10.8297 9.82973 12.7027 9.59294 15.0706L9.5 16H8.5L8.40706 15.0706C8.17027 12.7027 6.29734 10.8297 3.92939 10.5929L3 10.5V9.5L3.92939 9.40706C6.29734 9.17027 8.17027 7.29734 8.40706 4.92939Z"
				fill="currentColor"
			/>
		</svg>
	),
	media: ImageIcon,
	settings: GearIcon,
	billing: CreditCard,
	ellipsis: MoreVertical,
	add: PlusIcon,
	warning: ExclamationTriangleIcon,
	user: PersonIcon,
	arrowRight: ArrowRight,
	help: HelpCircle,
	sun: SunIcon,
	moon: MoonIcon,
	laptop: PersonIcon,
	message: ChatBubbleIcon,
	rocket: RocketIcon,
	edit: Pencil2Icon,
	link: Link2Icon,
	xCircle: CrossCircledIcon,
	pause: PauseIcon,
	play: PlayIcon,
	removeUser: MinusCircledIcon,
	rows: RowsIcon,
	gitHub: ({ ...props }: LucideProps) => (
		<svg
			aria-hidden="true"
			focusable="false"
			data-prefix="fab"
			data-icon="github"
			role="img"
			xmlns="http://www.w3.org/2000/svg"
			viewBox="0 0 496 512"
			{...props}
		>
			<path
				fill="currentColor"
				d="M165.9 397.4c0 2-2.3 3.6-5.2 3.6-3.3 .3-5.6-1.3-5.6-3.6 0-2 2.3-3.6 5.2-3.6 3-.3 5.6 1.3 5.6 3.6zm-31.1-4.5c-.7 2 1.3 4.3 4.3 4.9 2.6 1 5.6 0 6.2-2s-1.3-4.3-4.3-5.2c-2.6-.7-5.5 .3-6.2 2.3zm44.2-1.7c-2.9 .7-4.9 2.6-4.6 4.9 .3 2 2.9 3.3 5.9 2.6 2.9-.7 4.9-2.6 4.6-4.6-.3-1.9-3-3.2-5.9-2.9zM244.8 8C106.1 8 0 113.3 0 252c0 110.9 69.8 205.8 169.5 239.2 12.8 2.3 17.3-5.6 17.3-12.1 0-6.2-.3-40.4-.3-61.4 0 0-70 15-84.7-29.8 0 0-11.4-29.1-27.8-36.6 0 0-22.9-15.7 1.6-15.4 0 0 24.9 2 38.6 25.8 21.9 38.6 58.6 27.5 72.9 20.9 2.3-16 8.8-27.1 16-33.7-55.9-6.2-112.3-14.3-112.3-110.5 0-27.5 7.6-41.3 23.6-58.9-2.6-6.5-11.1-33.3 2.6-67.9 20.9-6.5 69 27 69 27 20-5.6 41.5-8.5 62.8-8.5s42.8 2.9 62.8 8.5c0 0 48.1-33.6 69-27 13.7 34.7 5.2 61.4 2.6 67.9 16 17.7 25.8 31.5 25.8 58.9 0 96.5-58.9 104.2-114.8 110.5 9.2 7.9 17 22.9 17 46.4 0 33.7-.3 75.4-.3 83.6 0 6.5 4.6 14.4 17.3 12.1C428.2 457.8 496 362.9 496 252 496 113.3 383.5 8 244.8 8zM97.2 352.9c-1.3 1-1 3.3 .7 5.2 1.6 1.6 3.9 2.3 5.2 1 1.3-1 1-3.3-.7-5.2-1.6-1.6-3.9-2.3-5.2-1zm-10.8-8.1c-.7 1.3 .3 2.9 2.3 3.9 1.6 1 3.6 .7 4.3-.7 .7-1.3-.3-2.9-2.3-3.9-2-.6-3.6-.3-4.3 .7zm32.4 35.6c-1.6 1.3-1 4.3 1.3 6.2 2.3 2.3 5.2 2.6 6.5 1 1.3-1.3 .7-4.3-1.3-6.2-2.2-2.3-5.2-2.6-6.5-1zm-11.4-14.7c-1.6 1-1.6 3.6 0 5.9 1.6 2.3 4.3 3.3 5.6 2.3 1.6-1.3 1.6-3.9 0-6.2-1.4-2.3-4-3.3-5.6-2z"
			></path>
		</svg>
	),
	twitter: TwitterLogoIcon,
	check: CheckIcon
}
