// Vite is kept minimal here because the app is a static SPA with a custom base path and test config.
// Public API: default Vite config export for dev server, build base path, and Vitest options.
// Constraints: configuration stays framework-agnostic and keeps tests deterministic in Node by default.
import { resolve } from 'node:path';
import { defineConfig } from 'vite';

export default defineConfig(({ command }) => ({
  // Dev keeps root (/) available for the status landing page, while build keeps GitHub Pages subpath assets.
  base: command === 'build' ? '/clonernews/' : '/',
  // Multi-page inputs keep a lightweight status landing at / and the SPA shell at /clonernews/.
  build: {
    rollupOptions: {
      input: {
        landing: resolve(__dirname, 'index.html'),
        app: resolve(__dirname, 'clonernews/index.html'),
      },
    },
  },
  // Vitest stays colocated with Vite so the same workspace config drives both build and test runs.
  test: {
    // Node is enough for these adapter and entity tests, while jsdom is reserved for DOM-specific cases.
    environment: 'node',
    globals: true,
    include: ['tests/**/*.test.js'],
    clearMocks: true,
  },
}));
