import { dummyEnv } from "../__fixtures__/dummy-environment";
import { environmentsSlice, fetchEnvironments } from "./environments";

describe("environmentsSlice reducer", () => {
  it("should handle initial state", () => {
    expect(
      environmentsSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toEqual({
      entities: {},
      ids: [],
    });
  });

  describe("fetchEnvironments", () => {
    it(`should handle ${fetchEnvironments.fulfilled.type}`, () => {
      expect(
        environmentsSlice.reducer(undefined, {
          type: fetchEnvironments.fulfilled.type,
          payload: [dummyEnv],
        })
      ).toEqual({
        entities: { [dummyEnv.id]: dummyEnv },
        ids: [dummyEnv.id],
      });
    });
  });
});
