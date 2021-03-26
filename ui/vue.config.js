const CompressionPlugin = require("compression-webpack-plugin");

module.exports = {
  publicPath: "./",
  assetsDir: "assets",
  productionSourceMap: false,
  devServer: {
    open: true,
    port: 9693
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
