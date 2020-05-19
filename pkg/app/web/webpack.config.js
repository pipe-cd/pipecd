/* eslint @typescript-eslint/no-var-requires: 0 */
"use strict";
const commonConfig = require("./webpack.common");
const mergeConfig = require("webpack-merge");

module.exports = (env, argv) => {
  return mergeConfig(commonConfig(env, argv), {
    mode: "production",
  });
};
