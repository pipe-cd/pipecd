import { loginSlice } from "./login";

describe("loginSlice reducer", () => {
  it("should handle initial state", () => {
    expect(
      loginSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toMatchInlineSnapshot(`
      Object {
        "projectName": null,
      }
    `);
  });
});
