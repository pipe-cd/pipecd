const ForkTsCheckerWebpackPlugin = require("fork-ts-checker-webpack-plugin");

module.exports = ({ config }) => {
  config.module.rules.push({
    test: /\.(ts|tsx)$/,
    use: [
      {
        loader: "ts-loader",
        options: {
          transpileOnly: true,
        },
      },
      require.resolve("react-docgen-typescript-loader"),
    ],
  });
  config.resolve.extensions.push(".ts", ".tsx");
  return config;
};
