/* eslint @typescript-eslint/no-var-requires: 0 */
"use strict";
const commonConfig = require("./webpack.common");
const { merge } = require("webpack-merge");
const path = require("path");

module.exports = (env) => {
  return merge(commonConfig(env), {
    resolve: {
      extensions: [".mjs", ".js", ".jsx"],
      alias: {
        pipe: path.resolve(env.bazelBinPath),
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
  });
};
