import { fileURLToPath, URL } from "node:url"
import os from "os"
import path from "path"
import fs from "fs"

import { defineConfig } from "vite"
import vue from "@vitejs/plugin-vue"
import vueDevTools from "vite-plugin-vue-devtools"

// element-plus auto import
import AutoImport from "unplugin-auto-import/vite"
import Components from "unplugin-vue-components/vite"
import { ElementPlusResolver } from "unplugin-vue-components/resolvers"

import packageJson from "./package.json"

// https://vite.dev/config/
// noinspection JSUnusedGlobalSymbols
export default defineConfig({
    define: {
        __PackageVersion__: JSON.stringify(packageJson.version),
        __IsDev__: process.env.NODE_ENV === "development", // 想要在router里使用，不将测试路由打包进生成文件，失败了
    },
    envPrefix: "Vite_",
    plugins: [
        vue(),
        vueDevTools(),

        AutoImport({ resolvers: [ ElementPlusResolver() ] }),
        Components({ resolvers: [ ElementPlusResolver() ] }),

        {
            name: "exclude-wasm-from-dist",
            apply: "build",
            closeBundle() {
                const removeFile = (fileName: string) => {
                    const file = path.resolve(__dirname, "dist/" + fileName)
                    if (fs.existsSync(file)) {
                        fs.unlinkSync(file)
                    }
                }

                removeFile("flip.wasm")
                removeFile("flip.html")
                removeFile("wasm_exec.js")
            }
        }
    ],
    resolve: { alias: { "@": fileURLToPath(new URL("./src", import.meta.url)) } },
    clearScreen: false,
    server: {
        host: getLocalIP(),
        port: 20319,
        open: "/#/test"
    }
})

export function getLocalIP(): string {
    const networks = os.networkInterfaces()
    for (const key in networks) {
        // @ts-ignore
        for (const ins of networks[key]) {
            if (ins.family === "IPv4" && !ins.internal) {
                return ins.address
            }
        }
    }

    return "127.0.0.1"
}
