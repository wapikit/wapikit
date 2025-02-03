import { AUTH_TOKEN_LS, getBackendUrl } from '~/constants'

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

	const response = await fetch(
		`${getBackendUrl()}${url}` + `${queryParam ? `?${queryParam}` : ''}`,
		{
			method,
			...(data ? { body: JSON.stringify(data) } : {}),
			headers: headers,
			credentials: 'include',
			mode: 'cors',
			cache: 'no-cache'
		}
	)

	// if (!response.ok) {
	// 	// Gracefully return an error object
	// 	const errorData = await response.json().catch(() => ({})) // Handle non-JSON error responses
	// 	return Promise.reject({
	// 		status: response.status,
	// 		statusText: response.statusText,
	// 		message: errorData.message || 'An error occurred'
	// 	})
	// }

	const responseData = await response.json()
	return responseData
}

export default customInstance
