import { commandsSlice } from "./commands";

describe("commandsSlice reducer", () => {
  it("should handle initial state", () => {
    expect(
      commandsSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toMatchInlineSnapshot(`
      Object {
        "entities": Object {},
        "ids": Array [],
      }
    `);
  });
});
