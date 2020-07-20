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
        "displayLength": 30,
        "entities": Object {},
        "ids": Array [],
        "loading": Object {},
        "loadingList": false,
      }
    `);
  });
});
