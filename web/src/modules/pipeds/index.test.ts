import { dummyPiped } from "~/__fixtures__/dummy-piped";
import { pipedsSlice, fetchPipeds } from "./";

const baseState = {
  entities: {},
  ids: [],
};

describe("pipedsSlice reducer", () => {
  it("should return the initial state", () => {
    expect(
      pipedsSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toEqual(baseState);
  });

  describe("fetchPipeds", () => {
    it(`should handle ${fetchPipeds.fulfilled.type}`, () => {
      expect(
        pipedsSlice.reducer(baseState, {
          type: fetchPipeds.fulfilled.type,
          payload: [dummyPiped],
        })
      ).toEqual({
        ...baseState,
        entities: { [dummyPiped.id]: dummyPiped },
        ids: [dummyPiped.id],
      });
    });
  });
});
