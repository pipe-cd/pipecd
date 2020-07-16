import { applicationFilterOptionsSlice } from "./application-filter-options";

describe("applicationFilterOptionsSlice reducer", () => {
  it("should handle initial state", () => {
    expect(
      applicationFilterOptionsSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toMatchInlineSnapshot(`
      Object {
        "enabled": Object {
          "value": true,
        },
        "envIdsList": Array [],
        "kindsList": Array [],
        "syncStatusesList": Array [],
      }
    `);
  });
});
