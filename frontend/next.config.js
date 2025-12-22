/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  swcMinify: true,
  env: {
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1',
    NEXT_PUBLIC_WS_URL: process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080/api/v1/ws',
  },
  images: {
    domains: ['localhost'],
  },
  webpack: (config, { isServer }) => {
    // undiciのプライベートフィールド構文の問題を回避
    if (!isServer) {
      config.resolve.fallback = {
        ...config.resolve.fallback,
        fs: false,
        net: false,
        tls: false,
        crypto: false,
      };
      // undiciをクライアントサイドでは使用しない
      config.resolve.alias = {
        ...config.resolve.alias,
        undici: false,
      };
    }
    return config;
  },
}

module.exports = nextConfig

