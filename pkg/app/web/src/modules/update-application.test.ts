import { updateApplicationSlice } from "./update-application";

describe("updateApplicationSlice reducer", () => {
  it("should handle initial state", () => {
    expect(
      updateApplicationSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toMatchInlineSnapshot();
  });
});