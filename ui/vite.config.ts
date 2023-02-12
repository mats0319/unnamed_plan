import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import viteCompression from "vite-plugin-compression"

import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'

import * as path from "path"

// reference: https://cn.vitejs.dev/config/
export default defineConfig({
  base: "./",
  envPrefix: "Vite_", // use: 'import.meta.env'
  plugins: [
    vue(),

    viteCompression(),

    AutoImport({
      resolvers: [
        ElementPlusResolver(),
      ]
    }),

    Components({
      resolvers: [
        ElementPlusResolver(),
      ]
    }),
  ],
  resolve: {
    alias: [
      {
        find: "@",
        replacement: path.resolve("./src")
      }
    ]
  },
  server: {
    open: true,
    host: "127.0.0.1",
    port: 9694,
  },
})
