/* eslint @typescript-eslint/no-var-requires: 0 */
"use strict";
const commonConfig = require("./gen_webpack.common");
const { merge } = require("webpack-merge");
const path = require("path");
const webpack = require("webpack");

// overridden by bazel
const version = "unknown_placeholder";

module.exports = (env) => {
  return merge(commonConfig(env), {
    resolve: {
      extensions: [".mjs", ".js", ".jsx"],
      alias: {
        pipecd: path.resolve(env.bazelBinPath),
        "~": path.resolve(env.bazelBinPath, "web/src"),
        "~~": path.resolve(env.bazelBinPath, "web"),
      },
    },
    module: {
      rules: [
        {
          type: "javascript/auto",
          test: /\.m?js$/,
          resolve: { fullySpecified: false },
          use: [],
        },
      ],
    },
    mode: "production",
    plugins: [
      new webpack.EnvironmentPlugin({
        STABLE_VERSION: version || null,
      }),
    ],
  });
};
