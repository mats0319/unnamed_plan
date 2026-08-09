interface ViteTypeOptions {
    // 将 ImportMetaEnv 的类型设为严格模式，不允许有未知的键值了。
    strictImportMetaEnv: unknown
}

interface ImportMetaEnv {
    DEV: boolean
    readonly Vite_axios_base_url: string
    readonly Vite_axios_flip_game_url: string
}

interface ImportMeta {
    readonly env: ImportMetaEnv
}
