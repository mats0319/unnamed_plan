declare global {
    const __PackageVersion__: string
    const __IsDev__: boolean
}

export {} // 加上这一行从全局声明，变成模块增强，不会污染整个
