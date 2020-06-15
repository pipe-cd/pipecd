import { activeStageSlice } from "./active-stage";

describe("activeStageSlice reducer", () => {
  it("should handle initial state", () => {
    expect(
      activeStageSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toMatchInlineSnapshot(`null`);
  });
});
