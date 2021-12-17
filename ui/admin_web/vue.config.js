const CompressionPlugin = require("compression-webpack-plugin");

module.exports = {
  publicPath: "./",
  outputDir: "dist_admin_web",
  assetsDir: "assets",
  productionSourceMap: false,
  pages: {
    index: {
      entry: "src/main.ts",
      template: "public/index.html",
      filename: "index.html"
    },
    unsupported: {
      entry: "src/views/unsupported/main.ts",
      template: "src/views/unsupported/unsupported.html",
      filename: "unsupported.html"
    }
  },
  css: {
    extract: false
  },
  devServer: {
    open: true,
    port: 19694
  },
  configureWebpack: {
    plugins: [
      new CompressionPlugin({
        // default algorithm: gzip
        test: new RegExp(/\.(js|css|svg)$/),
        threshold: 10240, // 10k
        minRatio: 0.8
      })
    ]
  }
}
