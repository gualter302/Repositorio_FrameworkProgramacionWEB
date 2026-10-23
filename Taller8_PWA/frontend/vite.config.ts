import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'
import { VitePWA } from 'vite-plugin-pwa'

export default defineConfig({
  plugins: [react(), tailwindcss(), VitePWA({
    registerType: 'prompt',
    includeAssets: ['icons/*.png', 'placeholder.svg'],
    manifest: {
      id: '/', name: 'MultiCatálogo', short_name: 'MultiCatálogo', lang: 'es',
      description: 'Tu catálogo y carrito, también sin conexión',
      start_url: '/catalogo', scope: '/', display: 'standalone',
      theme_color: '#020814', background_color: '#020814',
      icons: [
        { src: '/icons/icon-192.png', sizes: '192x192', type: 'image/png' },
        { src: '/icons/icon-512.png', sizes: '512x512', type: 'image/png' },
        { src: '/icons/maskable-512.png', sizes: '512x512', type: 'image/png', purpose: 'maskable' },
      ],
    },
    workbox: {
      globPatterns: ['**/*.{js,css,html,png,svg,webp,woff2}'],
      navigateFallback: '/index.html',
      navigateFallbackDenylist: [/^\/api(?:\/|$)/],
      cleanupOutdatedCaches: true,
      // Las respuestas de API se gestionan en el cliente: nunca se cachean datos privados.
      runtimeCaching: [{
        urlPattern: ({ request, url }) => request.destination === 'image' && url.hostname.endsWith('picsum.photos'),
        handler: 'CacheFirst',
        options: { cacheName: 'multicatalogo-images', expiration: { maxEntries: 40, maxAgeSeconds: 604800 }, cacheableResponse: { statuses: [0, 200] } },
      }],
    },
  })],
  server: { proxy: { '/api': 'http://127.0.0.1:3008' } },
  preview: { proxy: { '/api': 'http://127.0.0.1:3008' } },
})
