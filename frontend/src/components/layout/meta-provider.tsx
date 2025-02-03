'use client'

import { useEffect } from 'react'
import { useGetOrganizationTags, useGetUserFeatureFlags } from 'root/.generated'
import { useAuthState } from '~/hooks/use-auth-state'
import { useLayoutStore } from '~/store/layout.store'

const MetaProvider = ({ children }: { children: React.ReactNode }) => {
	const { authState } = useAuthState()

	const { data: featureFlags } = useGetUserFeatureFlags({
		query: {
			enabled: !!authState.isAuthenticated
		}
	})

	const { data: tags } = useGetOrganizationTags(
		{
			page: 1,
			per_page: 50
		},
		{
			query: {
				enabled: !!authState.isAuthenticated
			}
		}
	)

	const { writeProperty } = useLayoutStore()

	useEffect(() => {
		writeProperty({
			tags: tags?.tags || []
		})
	}, [tags?.tags, writeProperty])

	useEffect(() => {
		writeProperty({
			featureFlags: featureFlags?.featureFlags || null
		})
	}, [featureFlags?.featureFlags, writeProperty])

	return <>{children}</>
}

export default MetaProvider
