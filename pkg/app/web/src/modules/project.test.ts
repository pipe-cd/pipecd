import { projectSlice } from "./project";

describe("projectSlice reducer", () => {
  it("should handle initial state", () => {
    expect(
      projectSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toMatchInlineSnapshot(`
      Object {
        "desc": null,
        "github": null,
        "id": null,
        "isUpdatingGitHubSSO": false,
        "isUpdatingStaticAdmin": false,
        "staticAdminDisabled": false,
        "teams": null,
        "username": null,
      }
    `);
  });
});
