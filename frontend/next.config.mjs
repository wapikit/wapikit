/** @type {import('next').NextConfig} */
const nextConfig = {
	env: {
		NEXT_PUBLIC_IS_MANAGED_CLOUD_EDITION: process.env.NEXT_PUBLIC_IS_MANAGED_CLOUD_EDITION
	},
	output: 'export',
	images: {
		unoptimized: true,
		remotePatterns: [
			{
				hostname: 'unsplash.com'
			},
			{
				hostname: 'cdn.bfldr.com'
			}
		]
	}
}

export default nextConfig
