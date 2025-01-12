/** @type {import('next').NextConfig} */
const nextConfig = {
	env: {
		WAPIKIT_app__backend_url: process.env.WAPIKIT_app__backend_url
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
