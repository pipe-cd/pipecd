jest.spyOn(Date, "now").mockImplementation(() => 1);

import { insightSlice } from "./insight";

describe("insightSlice reducer", () => {
  it("should handle initial state", () => {
    expect(
      insightSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toMatchInlineSnapshot(`
      Object {
        "applicationId": "",
        "rangeFrom": 1,
        "rangeTo": 604800001,
        "step": 0,
      }
    `);
  });
});
