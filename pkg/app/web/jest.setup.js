"use strict";

let mockHistory = null;
jest.mock("./src/history.ts", () => ({
  __setMockHistory(his) {
    mockHistory = his;
  },
  get history() {
    return mockHistory;
  },
}));
