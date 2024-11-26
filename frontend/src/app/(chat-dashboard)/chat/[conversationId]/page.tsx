export async function generateStaticParams() {
	return [
		{
			userId: 'ok'
		}
	]
}

const UserChatPage = () => {
	// ! TODO:
	// ! get the conversation id
	// ! fetch the conversation and message in the reverse sort order with paginated
	// ! apply reach way point to fetch more messages

	return (
		<div>
			<h1>Chat</h1>
		</div>
	)
}

export default UserChatPage
