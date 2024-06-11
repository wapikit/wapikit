export async function generateStaticParams() {
	return [
		{
			userId: 'ok'
		}
	]
}

export const dynamicParams = true

const UserChatPage = () => {
	return (
		<div>
			<h1>Chat</h1>
		</div>
	)
}

export default UserChatPage
