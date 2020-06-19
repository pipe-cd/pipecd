import { pipedsSlice } from "./pipeds";

describe("pipedsSlice reducer", () => {
  it("should handle initial state", () => {
    expect(
      pipedsSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toMatchInlineSnapshot(`
      Object {
        "registeredPiped": null,
      }
    `);
  });
});
