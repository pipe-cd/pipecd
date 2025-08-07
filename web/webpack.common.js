/* eslint @typescript-eslint/no-var-requires: 0 */
"use strict";
const HtmlWebpackPlugin = require("html-webpack-plugin");
const path = require("path");
const webpack = require("webpack");
const CopyPlugin = require("copy-webpack-plugin");

module.exports = (env) => {
  return {
    output: {
      filename: "assets/[name].[contenthash:9].js",
      chunkFilename: "assets/[name].[contenthash:9].chunk.js",
      publicPath: "/",
      path: path.resolve(__dirname, "../.artifacts/web-static"),
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
          type: "asset/resource",
          generator: {
            filename: "assets/[name].[hash:8][ext]",
          },
        },
        {
          test: /\.css$/i,
          use: ["style-loader", "css-loader"],
        },
      ],
    },
    resolve: {
      fallback: {
        path: require.resolve("path-browserify"),
      },
    },
    plugins: [
      env.htmlTemplate &&
        new HtmlWebpackPlugin({
          filename: "index.html",
          template: env.htmlTemplate,
          favicon: path.join(__dirname, "assets/favicon.ico"),
        }),
      process.env.ENABLE_MOCK &&
        new CopyPlugin({
          patterns: [
            {
              from: path.join(__dirname, "public/mockServiceWorker.js"),
            },
          ],
        }),
      new webpack.EnvironmentPlugin({
        ENABLE_MOCK: process.env.ENABLE_MOCK || null,
      }),
      new webpack.ProvidePlugin({
        process: "process/browser",
      }),
    ].filter(Boolean),
  };
};
