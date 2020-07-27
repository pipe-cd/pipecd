import { meSlice } from "./me";

describe("meSlice reducer", () => {
  it("should handle initial state", () => {
    expect(
      meSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toMatchInlineSnapshot(`null`);
  });
});
