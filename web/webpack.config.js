/* eslint @typescript-eslint/no-var-requires: 0 */
"use strict";
require("dotenv").config();
const path = require("path");
const webpack = require("webpack");
const { merge } = require("webpack-merge");
const webpackBaseConfig = require("./webpack.common");

module.exports = (env) => {
  return merge(webpackBaseConfig(env), {
    mode: "production",
    entry: {
      index: "./src/index.tsx",
    },
    resolve: {
      //extensions: [".mjs", ".js", ".jsx"],
      extensions: [".mjs", ".ts", ".tsx", ".js"],
      alias: {
        pipecd: path.resolve(__dirname, ".."),
        "~": path.resolve(__dirname, "src"),
        "~~": path.resolve(__dirname),
      },
      modules: [path.resolve(__dirname, "node_modules"), "node_modules"],
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
          use: [],
        },
      ],
    },
    plugins: [
      new webpack.EnvironmentPlugin({
        STABLE_VERSION: process.env.STABLE_VERSION || null,
      }),
    ],
  });
};
