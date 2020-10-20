import { sealedSecretSlice } from "./sealed-secret";

describe("sealedSecretSlice reducer", () => {
  it("should handle initial state", () => {
    expect(
      sealedSecretSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toMatchInlineSnapshot(`
      Object {
        "data": null,
        "isLoading": false,
      }
    `);
  });
});
