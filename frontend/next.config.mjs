/** @type {import('next').NextConfig} */
const nextConfig = {
	env: {
		BACKEND_URL: process.env.BACKEND_URL
	},
	output: 'export'
}

export default nextConfig
