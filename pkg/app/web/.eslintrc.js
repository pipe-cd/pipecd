module.exports = {
  extends: [
    "eslint:recommended",
    "plugin:@typescript-eslint/recommended",
    "plugin:@typescript-eslint/eslint-recommended",
    "plugin:prettier/recommended",
    "prettier/@typescript-eslint",
    "plugin:react/recommended",
  ],
  plugins: ["@typescript-eslint", "react-hooks", "jest"],
  parser: "@typescript-eslint/parser",
  env: { browser: true, node: true, es6: true, "jest/globals": true },
  rules: {
    "react-hooks/rules-of-hooks": "error",
    "react-hooks/exhaustive-deps": "warn",
    "react/prop-types": "off",
    "@typescript-eslint/explicit-function-return-type": [
      "error",
      {
        allowExpressions: true,
      },
    ],
  },
  settings: {
    react: {
      version: "detect",
    },
  },
};
