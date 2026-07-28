import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  // relative base so the build works when served from a GitHub Pages
  // project subpath (https://<user>.github.io/<repo>/) without needing to
  // know the repo name ahead of time
  base: './',
  plugins: [react()],
})
