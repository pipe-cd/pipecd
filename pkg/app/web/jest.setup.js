"use strict";

Object.defineProperty(navigator, "clipboard", {
  value: {
    writeText: (content) => {
      return Promise.resolve(content);
    },
  },
});

let mockHistory = null;
jest.mock("./src/history.ts", () => ({
  __setMockHistory(his) {
    mockHistory = his;
  },
  get history() {
    return mockHistory;
  },
}));
