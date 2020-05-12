/* eslint @typescript-eslint/no-var-requires: 0 */
"use strict";
const path = require("path");
const HtmlWebpackPlugin = require("html-webpack-plugin");

module.exports = (_, argv) => {
  return {
    mode: "production",
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
      },
    },
    resolve: {
      extensions: [".mjs", ".js", ".jsx"],
      alias: {
        pipe: path.resolve(argv.testPath),
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

              return "[name].[hash:8].[ext]";
            },
          },
        },
        {
          type: "javascript/auto",
          test: /\.mjs$/,
          use: [],
        },
      ],
    },
    plugins: [
      argv.htmlTemplate &&
        new HtmlWebpackPlugin({
          filename: "index.html",
          template: argv.htmlTemplate,
        }),
    ].filter(Boolean),
  };
};
