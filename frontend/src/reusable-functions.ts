import { ulid } from 'ulid'

export function generateUniqueId() {
	return ulid()
}

export function getWebsocketUrl(token: string) {
	return process.env.NODE_ENV === 'development' ? `ws://localhost:3001?token=${token}` : ``
}
