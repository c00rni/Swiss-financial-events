import path from "path"
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

const root = path.resolve(__dirname, "./src")
const outDir = path.resolve(__dirname, "./_dist/")

// https://vitejs.dev/config/
export default defineConfig({
    root,
    plugins: [react()],
    resolve: {
        alias: {
            "@": path.resolve(__dirname, "./src"),
        },
    },
    build: {
        outDir,
        emptyOutDir: true,
        rollupOptions: {
            input: {
                main: path.resolve(root, "index.html"),
            }
        }
    }
})

