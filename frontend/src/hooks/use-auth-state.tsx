'use client'

import { AUTH_TOKEN_LS } from '~/constants'
import { useEffect, useState } from 'react'
import { z } from 'zod'
import { UserPermissionLevelEnum } from 'root/.generated'
import { decode } from 'jsonwebtoken'
import { UserTokenPayloadSchema } from '~/schema'

const AuthStateSchemaType = z
	.object({
		isAuthenticated: z.literal(true),
		data: z.object({
			user: z.object({
				uniqueId: z.string(),
				email: z.string(),
				role: z.nativeEnum(UserPermissionLevelEnum).or(z.string().nullish()),
				username: z.string(),
				organizationId: z.string().nullish(),
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
	const [authState, setAuthState] = useState<z.infer<typeof AuthStateSchemaType>>({
		isAuthenticated: null
	})

	useEffect(() => {
		const authToken = localStorage.getItem(AUTH_TOKEN_LS)

		if (authToken) {
			// decode the json web token here
			const payload = decode(authToken)
			const parsedPayload = UserTokenPayloadSchema.safeParse(payload)

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
				setAuthState(() => ({ isAuthenticated: false }))
			}
		} else {
			setAuthState(() => ({ isAuthenticated: false }))
		}
	}, [])

	return {
		authState
	}
}
