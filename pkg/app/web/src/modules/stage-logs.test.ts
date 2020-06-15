import { stageLogsSlice } from "./stage-logs";

describe("stageLogsSlice reducer", () => {
  it("should handle initial state", () => {
    expect(
      stageLogsSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toMatchInlineSnapshot(`Object {}`);
  });
});
