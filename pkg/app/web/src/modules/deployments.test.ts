import { deploymentsSlice } from "./deployments";

describe("deploymentsSlice reducer", () => {
  it("should handle initial state", () => {
    expect(
      deploymentsSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toMatchInlineSnapshot(`
      Object {
        "canceling": Object {},
        "entities": Object {},
        "hasMore": true,
        "ids": Array [],
        "isLoadingItems": false,
        "isLoadingMoreItems": false,
        "loading": Object {},
      }
    `);
  });
});
