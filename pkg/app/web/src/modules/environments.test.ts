import { environmentsSlice } from "./environments";

describe("environmentsSlice reducer", () => {
  it("should handle initial state", () => {
    expect(
      environmentsSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toMatchInlineSnapshot(`
      Object {
        "entities": Object {},
        "ids": Array [],
      }
    `);
  });
});
