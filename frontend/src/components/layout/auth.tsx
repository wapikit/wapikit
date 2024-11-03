'use client'

import { usePathname, useRouter } from 'next/navigation'
import React, { useEffect } from 'react'
import { useAuthState } from '~/hooks/use-auth-state'
import LoadingSpinner from '../loader'
import { useGetUser } from 'root/.generated'
import { useLayoutStore } from '~/store/layout.store'

const AuthProvisioner: React.FC<{ children: React.ReactNode }> = ({ children }) => {
	const { authState } = useAuthState()
	const router = useRouter()
	const pathname = usePathname()

	const { writeProperty } = useLayoutStore()

	useEffect(() => {
		console.log('authState.isAuthenticated', authState.isAuthenticated)

		if (pathname === '/signin') {
			return
		} else {
			if (authState.isAuthenticated === false) {
				router.push('/signin')
			} else {
				// either auth is loading or user is authenticated
			}
		}
	}, [authState.isAuthenticated, pathname, router])

	const { data: userData } = useGetUser({
		query: {
			enabled: !!authState.isAuthenticated
		}
	})

	useEffect(() => {
		if (!authState.isAuthenticated || !userData) {
			return
		}

		writeProperty({
			user: userData
		})
	}, [userData, authState.isAuthenticated, writeProperty])

	if (
		typeof authState.isAuthenticated !== 'boolean' &&
		!authState.isAuthenticated &&
		pathname !== '/'
	) {
		return (
			<div className="flex h-full w-full items-center justify-center">
				<LoadingSpinner />
			</div>
		)
	} else {
		return <>{children}</>
	}
}

export default AuthProvisioner
