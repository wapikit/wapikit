import { AUTH_TOKEN_LS, BACKEND_URL } from '~/constants'

export const customInstance = async <T>({
	url,
	method,
	params,
	data
}: {
	url: string
	method: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH'
	params?: any
	data?: any
	responseType?: string
	signal?: AbortSignal
	headers?: Record<string, string>
}): Promise<T> => {
	const authToken = localStorage.getItem(AUTH_TOKEN_LS)
	const headers = new Headers()
	headers.set('Content-Type', 'application/json')
	headers.set('Accept', 'application/json')

	if (authToken) {
		console.log({ authToken })
		headers.set('x-access-token', authToken)
	}

	const response = await fetch(`${BACKEND_URL}${url}` + `?` + new URLSearchParams(params), {
		method,
		...(data ? { body: JSON.stringify(data) } : {}),
		headers: headers,
		credentials: 'include',
		mode: 'cors',
		cache: 'no-cache'
	})

	console.log({ response })

	return response.json()
}

export default customInstance
