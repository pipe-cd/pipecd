import { applicationCountsSlice } from "./application-counts";

describe("applicationCountsSlice reducer", () => {
  it("should return the initial state", () => {
    expect(
      applicationCountsSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toMatchInlineSnapshot(`
      Object {
        "counts": Object {
          "CLOUDRUN": Object {
            "DISABLED": 0,
            "ENABLED": 0,
          },
          "CROSSPLANE": Object {
            "DISABLED": 0,
            "ENABLED": 0,
          },
          "ECS": Object {
            "DISABLED": 0,
            "ENABLED": 0,
          },
          "KUBERNETES": Object {
            "DISABLED": 0,
            "ENABLED": 0,
          },
          "LAMBDA": Object {
            "DISABLED": 0,
            "ENABLED": 0,
          },
          "TERRAFORM": Object {
            "DISABLED": 0,
            "ENABLED": 0,
          },
        },
        "updatedAt": 0,
      }
    `);
  });
});
