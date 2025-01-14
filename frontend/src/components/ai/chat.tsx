import { useChat } from 'ai/react'

export default function ChatComponent() {
	const { messages, input, handleInputChange, handleSubmit } = useChat({
		api: 'http://localhost:8080/v1/ai/chat'
	})

	return (
		<div>
			{messages.map((m, i) => (
				<div key={i}>
					<strong>{m.role}: </strong>
					{m.content}
				</div>
			))}
			<form onSubmit={handleSubmit}>
				<input value={input} onChange={handleInputChange} />
				<button type="submit">Send</button>
			</form>
		</div>
	)
}
