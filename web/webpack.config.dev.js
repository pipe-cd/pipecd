/* eslint @typescript-eslint/no-var-requires: 0 */
"use strict";
require("dotenv").config();
const path = require("path");
const webpack = require("webpack");
const ForkTsCheckerWebpackPlugin = require("fork-ts-checker-webpack-plugin");
const { merge } = require("webpack-merge");
const webpackBaseConfig = require("./webpack.common");

module.exports = (env) =>
  merge(webpackBaseConfig(env), {
    mode: process.env.NODE_ENV === "production" ? "production" : "development",
    devtool: "inline-source-map",
    entry: {
      index: "./src/index.tsx",
    },
    resolve: {
      extensions: [".mjs", ".ts", ".tsx", ".js"],
      alias: {
        pipecd: path.resolve(__dirname, ".."),
        "~": path.resolve(__dirname, "src"),
        "~~": path.resolve(__dirname),
      },
      modules: [path.resolve(__dirname, "node_modules"), "node_modules"],
    },
    devServer: {
      static: [path.join(__dirname, "dist"), path.join(__dirname, "public")],
      compress: true,
      port: 9090,
      historyApiFallback: true,
      allowedHosts: "all",
      proxy: [
        {
          context: ["/api"],
          changeOrigin: true,
          target: process.env.API_ENDPOINT,
          pathRewrite: { "^/api": "" },
          withCredentials: true,
          headers: {
            Cookie: process.env.API_COOKIE,
          },
        },
      ],
    },
    module: {
      rules: [
        {
          test: /\.tsx?$/,
          loader: "ts-loader",
          options: {
            transpileOnly: true,
          },
        },
        {
          type: "javascript/auto",
          test: /\.m?js$/,
          resolve: { fullySpecified: false },
        },
      ],
    },
    plugins: [
      new ForkTsCheckerWebpackPlugin(),
      new webpack.EnvironmentPlugin({
        API_ENDPOINT: process.env.API_ENDPOINT || null,
        PIPECD_VERSION: process.env.PIPECD_VERSION || null,
      }),
    ],
  });
