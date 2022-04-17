const CompressionPlugin = require("compression-webpack-plugin");

module.exports = {
  publicPath: "./",
  outputDir: "dist_public_mobile",
  assetsDir: "assets",
  productionSourceMap: false,
  pages: {
    index: {
      entry: "src/main.ts",
      template: "public/index.html",
      filename: "index.html"
    }
  },
  css: {
    extract: false
  },
  devServer: {
    // open: true,
    port: 9695
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
