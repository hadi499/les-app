import { sveltekit } from "@sveltejs/kit/vite";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [
    sveltekit(),
    tailwindcss(),
  ],
  server: {
    host: true,
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
      "^/me$": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
      "^/uploads": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
});
