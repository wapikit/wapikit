/** @type {import('next').NextConfig} */
const nextConfig = {
	env: {
		BACKEND_URL: process.env.BACKEND_URL
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
