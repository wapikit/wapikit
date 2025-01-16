'use client'

import { useEffect } from 'react'
import { useGetUserFeatureFlags } from 'root/.generated'
import { useLayoutStore } from '~/store/layout.store'

const FeatureFlagProvider = ({ children }: { children: React.ReactNode }) => {
	const { data: featureFlags } = useGetUserFeatureFlags()

	const { writeProperty } = useLayoutStore()

	useEffect(() => {
		writeProperty({
			featureFlags: featureFlags?.featureFlags || null
		})
	}, [featureFlags?.featureFlags, writeProperty])

	return <>{children}</>
}

export default FeatureFlagProvider
