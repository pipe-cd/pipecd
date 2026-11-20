module.exports = {
  roots: ["<rootDir>/src"],
  transform: {
    "^.+\\.[jt]sx?$": [
      "ts-jest",
      {
        isolatedModules: true,
      },
    ],
    "\\.(jpg|jpeg|png|gif|eot|otf|webp|svg|ttf|woff|woff2|mp4|webm|wav|mp3|m4a|aac|oga|ico)$":
      "<rootDir>/file-transformer.js",
  },
  // query-string and some of its dependencies are published as ESM only.
  transformIgnorePatterns: [
    "/node_modules/(?!(query-string|filter-obj|split-on-first|decode-uri-component)/)",
  ],
  moduleNameMapper: {
    "^pipecd/(.*)$": "<rootDir>/../$1",
    "^~/(.*)$": "<rootDir>/src/$1",
    "^~~/(.*)$": "<rootDir>/$1",
    "\\.(css|less|scss|sass)$": "<rootDir>/test-utils/styleMock.ts",
  },
  moduleDirectories: ["<rootDir>/node_modules", "__fixtures__"],
  coveragePathIgnorePatterns: ["/node_modules/", ".test.ts", ".d.ts"],
  testEnvironment: "./custom-jsdom",
  clearMocks: true,
  setupFiles: ["./jest.setup.js"],
  setupFilesAfterEnv: ["./jest.after-env.ts"],
  coverageReporters: ["lcovonly", "text-summary", "html"],
  maxWorkers: 4,
  testTimeout: 10000,
};
