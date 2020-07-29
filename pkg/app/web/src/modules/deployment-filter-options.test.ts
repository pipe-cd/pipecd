import { deploymentFilterOptionsSlice } from "./deployment-filter-options";

describe("deploymentFilterOptionsSlice reducer", () => {
  it("should handle initial state", () => {
    expect(
      deploymentFilterOptionsSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toMatchInlineSnapshot(`
      Object {
        "applicationIds": Array [],
        "envIds": Array [],
        "kinds": Array [],
        "statuses": Array [],
      }
    `);
  });
});
