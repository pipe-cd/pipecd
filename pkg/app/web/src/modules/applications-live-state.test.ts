import { applicationLiveStateSlice } from "./applications-live-state";

describe("applicationLiveStateSlice reducer", () => {
  it("should handle initial state", () => {
    expect(
      applicationLiveStateSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toMatchInlineSnapshot(`
      Object {
        "entities": Object {},
        "hasError": Object {},
        "ids": Array [],
      }
    `);
  });
});
