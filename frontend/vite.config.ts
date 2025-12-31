import tailwindcss from '@tailwindcss/vite';
import {defineConfig} from 'vite';
import {svelte} from '@sveltejs/vite-plugin-svelte';
import path from "path"

// https://vite.dev/config/
export default defineConfig({
    plugins: [tailwindcss(), svelte()], resolve: {
        alias: {
            $lib: path.resolve("./src/lib"),
            $backend: path.resolve("./wailsjs/go/cmd/App.js"),
            $runtime: path.resolve("./wailsjs/runtime/runtime.js")
        },
    },
});
