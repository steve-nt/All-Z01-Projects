import { defineConfig } from "vite";
import path from "node:path";

export default defineConfig({
  // Keep the runnable app under apps/, but allow running from repo root.
  root: "apps/todomvc",
  resolve: {
    alias: {
      // Consumer apps import the framework as "mini-framework"
      "mini-framework": path.resolve(
        process.cwd(),
        "packages/mini-framework/src/index.js",
      ),
    },
  },
  build: {
    // Output at repo root for convenience
    outDir: path.resolve(process.cwd(), "dist"),
    emptyOutDir: true,
  },
});

