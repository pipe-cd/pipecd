import { deploymentFrequencySlice } from "./deployment-frequency";

describe("deploymentFrequencySlice reducer", () => {
  it("should handle initial state", () => {
    expect(
      deploymentFrequencySlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toMatchInlineSnapshot(`
      Object {
        "data": Array [],
        "status": "idle",
      }
    `);
  });
});
