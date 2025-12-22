import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  // Load env file based on `mode` in the current working directory.
  const env = loadEnv(mode, process.cwd(), '')
  
  // Determine backend URL for proxy
  // If VITE_API_URL is relative (starts with /), use default backend URL
  // Otherwise, extract backend URL from VITE_API_URL
  const apiUrl = env.VITE_API_URL || '/api/v1'
  const backendUrl = apiUrl.startsWith('/') 
    ? 'http://localhost:8080'  // Changé de 8081 à 8082
    : apiUrl.replace(/\/api\/v1$/, '')
  
  return {
    plugins: [react()],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
    },
    server: {
      port: 3000,
      proxy: {
        '/api': {
          target: backendUrl,
          changeOrigin: true,
        },
      },
    },
    test: {
      globals: true,
      environment: 'jsdom',
      setupFiles: './src/test/setup.ts',
    },
  }
})
