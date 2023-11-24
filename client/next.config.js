/** @type {import('next').NextConfig} */

const rewrites =
  process.env.NODE_ENV === 'development'
    ? [
        {
          source: '/api/:path*',
          destination: 'http://localhost:80/:path*',
        },
      ]
    : [];

const nextConfig = {
  async rewrites() {
    return rewrites;
  },
};

module.exports = nextConfig;
