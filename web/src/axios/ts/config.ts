// Generate File, Should Not Edit.
// Author : mario. github.com/mats0319
// Code   : github.com/mats0319/study/go/gocts
// Version: gocts v0.2.5

import axios, { AxiosInstance } from "axios"

export const axiosWrapper: AxiosInstance = axios.create({
    baseURL: getBaseUrl(),
    timeout: 3000,
})

// getBaseUrl according to dev/product env & local ipv4 addr, set 'baseURL'
function getBaseUrl(): string {
    const url = import.meta.env.Vite_axios_base_url
    const localIP = window.location.hostname

    return import.meta.env.DEV ? url.replace("127.0.0.1", localIP) : url
}

// generate '.env.development' file in root path,
// with content 'Vite_axios_base_url = "http://127.0.0.1:10319/api"'
//
// generate '.env.production' file in root path,
// with content 'Vite_axios_base_url = "https://xxx.xxx"'
//
// modify 'vite.config.ts', add config:
// envPrefix: "Vite_"
