import { toastsSlice } from "./toasts";

describe("toastsSlice reducer", () => {
  it("should handle initial state", () => {
    expect(
      toastsSlice.reducer(undefined, {
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
