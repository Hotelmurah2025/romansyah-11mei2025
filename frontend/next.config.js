/** @type {import('next').NextConfig} */
module.exports = {
  reactStrictMode: true,
  output: 'export',
  distDir: 'dist',
  images: {
    unoptimized: true,
  },
  trailingSlash: false,
  assetPrefix: '',
  basePath: '',
  exportPathMap: async function() {
    return {
      '/': { page: '/' },
      '/login': { page: '/login' },
      '/register': { page: '/register' },
      '/dashboard': { page: '/dashboard' },
      '/hotels': { page: '/hotels' },
      '/hotels/new': { page: '/hotels/new' },
      '/rooms': { page: '/rooms' },
      '/bookings': { page: '/bookings' },
      '/admin/dashboard': { page: '/admin/dashboard' },
    };
  },
}
