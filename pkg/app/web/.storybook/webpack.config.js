const path = require("path");

module.exports = ({ config }) => {
  config.module.rules.push({
    type: "javascript/auto",
    test: /\.mjs$/,
    use: [],
  });
  config.resolve.extensions.push(".mjs", ".ts", ".tsx", ".js");
  config.resolve.modules.push(
    path.resolve(__dirname, "../node_modules"),
    "node_modules"
  );
  config.resolve.alias = {
    pipe: path.resolve(__dirname, "../../../../bazel-bin/"),
  };

  return config;
};
