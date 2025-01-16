/** @type {import('next').NextConfig} */
const nextConfig = {
	env: {},
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
