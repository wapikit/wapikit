import { BACKEND_URL } from '~/constants'

export const customInstance = async <T>({
	url,
	method,
	params,
	data,
}: {
	url: string
	method: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH'
	params?: any
	data?: any
	responseType?: string,
	signal?: AbortSignal
	headers?: Record<string, string>
}): Promise<T> => {
	// ! TODO: fetch the auth token here

	const authToken = localStorage.getItem('authToken')
	const headers = new Headers()
	headers.set('Content-Type', 'application/json')
	headers.set('Accept', 'application/json')

	if (authToken) {
		headers.set('x-access-token', authToken)
	}

	const response = await fetch(`${BACKEND_URL}${url}` + new URLSearchParams(params), {
		method,
		...(data ? { body: JSON.stringify(data) } : {}),
		headers: headers
	})

	return response.json()
}

export default customInstance
