'use client'

import { usePathname, useRouter } from 'next/navigation'
import React, { useEffect } from 'react'
import { useAuthState } from '~/hooks/use-auth-state'
import LoadingSpinner from '../loader'

const AuthProvisioner: React.FC<{ children: React.ReactNode }> = ({ children }) => {
	const { authState } = useAuthState()
	const router = useRouter()
	const pathname = usePathname()

	useEffect(() => {
		if (pathname === '/') {
			return
		} else {
			if (authState.isAuthenticated === false) {
				router.push('/')
			} else {
				// either auth is loading or user is authenticated
			}
		}
	}, [authState.isAuthenticated, pathname, router])

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
