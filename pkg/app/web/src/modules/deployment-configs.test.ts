import { deploymentConfigsSlice } from "./deployment-configs";

describe("deploymentConfigsSlice reducer", () => {
  it("should handle initial state", () => {
    expect(
      deploymentConfigsSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toMatchInlineSnapshot(`
      Object {
        "targetApplicationId": null,
        "templates": Object {},
      }
    `);
  });
});
