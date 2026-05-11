import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

export default defineConfig({
  // Part 1: enable Vue SFC support
  plugins: [vue()],
  server: {
    proxy: {
      // Part 1: proxy API/auth routes to Go backend (avoid CORS in dev)
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true
      },
      // Feed post images (backend serves uploads at /frontend/uploads)
      "/frontend/uploads": {
        target: "http://localhost:8080",
        changeOrigin: true
      },
      "/login": {
        target: "http://localhost:8080",
        changeOrigin: true
      },
      "/register": {
        target: "http://localhost:8080",
        changeOrigin: true
      },
      "/logout": {
        target: "http://localhost:8080",
        changeOrigin: true
      },
      "/health": {
        target: "http://localhost:8080",
        changeOrigin: true
      },
      // Part 1: WebSocket proxy for future realtime features
      "/ws": {
        target: "ws://localhost:8080",
        ws: true
      }
    }
  }
});
