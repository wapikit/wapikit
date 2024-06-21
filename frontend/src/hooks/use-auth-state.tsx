'use client'

import { AUTH_TOKEN_LS } from '~/constants'
import { useLocalStorage } from './use-local-storage'
import { useEffect, useState } from 'react'
import { z } from 'zod'
import { UserRoleEnum } from 'root/.generated'
import { decode } from 'jsonwebtoken'
import { UserTokenPayloadSchema } from '~/schema'

const AuthStateSchemaType = z
	.object({
		isAuthenticated: z.literal(true),
		data: z.object({
			user: z.object({
				uniqueId: z.string(),
				email: z.string(),
				role: z.nativeEnum(UserRoleEnum),
				username: z.string(),
				organizationId: z.string(),
				name: z.string()
			}),
			token: z.string()
		})
	})
	.or(
		z.object({
			isAuthenticated: z.literal(false).or(z.null())
		})
	)

export const useAuthState = () => {
	const [authToken] = useLocalStorage(AUTH_TOKEN_LS, '')

	const [authState, setAuthState] = useState<z.infer<typeof AuthStateSchemaType>>({
		isAuthenticated: null
	})

	useEffect(() => {
		console.log({ authToken })
		if (authToken) {
			// decode the json web token here
			const payload = decode(authToken)
			console.log({ payload })
			const parsedPayload = UserTokenPayloadSchema.safeParse(payload)
			console.log({ parsedPayload: parsedPayload.error?.errors })
			if (parsedPayload.success) {
				setAuthState(() => ({
					isAuthenticated: true,
					data: {
						token: authToken,
						user: {
							email: parsedPayload.data.email,
							uniqueId: parsedPayload.data.unique_id,
							username: parsedPayload.data.username,
							organizationId: parsedPayload.data.organization_id,
							role: parsedPayload.data.role,
							name: parsedPayload.data.name
						}
					}
				}))
			} else {
				// set auth to no
				setAuthState(() => ({ isAuthenticated: false }))
			}
		}
	}, [authToken])

	return {
		authState,
		authToken
	}
}
