import adapter from '@sveltejs/adapter-node';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
compilerOptions: {
  runes: true, // Will test if this breaks anything. Wait, svelte.config.js doesn't allow functions? Actually it does.
},
preprocess: vitePreprocess(),
kit: {
adapter: adapter()
}
};

export default config;
