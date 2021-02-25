/* eslint @typescript-eslint/no-var-requires: 0 */
"use strict";
const HtmlWebpackPlugin = require("html-webpack-plugin");
const path = require("path");
const webpack = require("webpack");

module.exports = (_, argv) => {
  return {
    output: {
      filename: "assets/[name].[contenthash:9].js",
      chunkFilename: "assets/[name].[contenthash:9].chunk.js",
      publicPath: "/",
    },
    optimization: {
      splitChunks: {
        chunks: "all",
        maxInitialRequests: 20, // for HTTP2
        maxAsyncRequests: 20, // for HTTP2
        cacheGroups: {
          service: {
            test: /[\\/]service/,
          },
        },
      },
    },
    module: {
      rules: [
        {
          test: /\.(png|svg|jpg|gif|ico)$/,
          loader: "file-loader",
          options: {
            name(file) {
              if (file.includes("/favicon.ico")) {
                return "[name].[ext]";
              }

              return "assets/[name].[hash:8].[ext]";
            },
          },
        },
      ],
    },
    plugins: [
      argv.htmlTemplate &&
        new HtmlWebpackPlugin({
          filename: "index.html",
          template: argv.htmlTemplate,
          favicon: path.join(__dirname, "assets/favicon.ico"),
        }),
      new webpack.EnvironmentPlugin(["NODE_ENV", "ENABLE_MOCK"]),
    ].filter(Boolean),
  };
};
