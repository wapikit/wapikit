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
		headers.set('x-access-token', authToken)
	}

	const queryParam = new URLSearchParams(params).toString()

	const response = await fetch(`${BACKEND_URL}${url}` + `${queryParam ? `?${queryParam}` : ''}`, {
		method,
		...(data ? { body: JSON.stringify(data) } : {}),
		headers: headers,
		credentials: 'include',
		mode: 'cors',
		cache: 'no-cache'
	})

	const responseData = await response.json()

	return responseData
}

export default customInstance
