'use client'

import { useEffect } from 'react'
import { useGetUserFeatureFlags } from 'root/.generated'
import { useAuthState } from '~/hooks/use-auth-state'
import { useLayoutStore } from '~/store/layout.store'

const FeatureFlagProvider = ({ children }: { children: React.ReactNode }) => {
	const { authState } = useAuthState()

	const { data: featureFlags } = useGetUserFeatureFlags({
		query: {
			enabled: !!authState.isAuthenticated
		}
	})

	const { writeProperty } = useLayoutStore()

	useEffect(() => {
		writeProperty({
			featureFlags: featureFlags?.featureFlags || null
		})
	}, [featureFlags?.featureFlags, writeProperty])

	return <>{children}</>
}

export default FeatureFlagProvider
