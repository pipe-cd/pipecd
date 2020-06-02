/* eslint @typescript-eslint/no-var-requires: 0 */
"use strict";
const commonConfig = require("./webpack.common");
const mergeConfig = require("webpack-merge");
const path = require("path");

module.exports = (env, argv) => {
  return mergeConfig(commonConfig(env, argv), {
    resolve: {
      extensions: [".mjs", ".js", ".jsx"],
      alias: {
        pipe: path.resolve(argv.bazelBinPath),
      },
    },
    module: {
      rules: [
        {
          type: "javascript/auto",
          test: /\.mjs$/,
          use: [],
        },
      ],
    },
    mode: "production",
  });
};
