import { useEffect, useState } from 'react'

export const useLocalStorage = <T>(key: string, initialValue: T): [T, (value: T) => void] => {
	const [storedValue, setStoredValue] = useState(initialValue)

	useEffect(() => {
		// Retrieve from localStorage
		const item = window.localStorage.getItem(key)
		if (item) {
			if (typeof item === 'string') {
				setStoredValue(item as T)
			} else {
				setStoredValue(JSON.parse(item))
			}
		}
	}, [key])

	const setValue = (value: T) => {
		// Save state
		setStoredValue(value)
		// Save to localStorage

		window.localStorage.setItem(key, typeof value === 'string' ? value : JSON.stringify(value))
	}
	return [storedValue, setValue]
}
