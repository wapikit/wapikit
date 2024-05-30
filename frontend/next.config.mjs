/** @type {import('next').NextConfig} */
const nextConfig = {
        env: {
            NODE_ENV: process.env.NODE_ENV,
            BACKEND_URL: process.env.BACKEND_URL,
        }
};

export default nextConfig;
