import { dummyPiped } from "~/__fixtures__/dummy-piped";
import {
  pipedsSlice,
  clearRegisteredPipedInfo,
  addPiped,
  fetchPipeds,
  addNewPipedKey,
  editPiped,
} from "./";

const baseState = {
  entities: {},
  ids: [],
  registeredPiped: null,
  updating: false,
  releasedVersions: [],
  breakingChangesNote: "",
};

describe("pipedsSlice reducer", () => {
  it("should return the initial state", () => {
    expect(
      pipedsSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toEqual(baseState);
  });

  it(`should handle ${clearRegisteredPipedInfo.type}`, () => {
    expect(
      pipedsSlice.reducer(
        {
          ...baseState,
          registeredPiped: {
            id: "piped-1",
            key: "piped-key",
            isNewKey: false,
          },
        },
        {
          type: clearRegisteredPipedInfo.type,
        }
      )
    ).toEqual(baseState);
  });

  describe("addPiped", () => {
    it(`should handle ${addPiped.fulfilled.type}`, () => {
      expect(
        pipedsSlice.reducer(baseState, {
          type: addPiped.fulfilled.type,
          payload: {
            id: "piped-1",
            key: "piped-key",
            isNewKey: false,
          },
        })
      ).toEqual({
        ...baseState,
        registeredPiped: {
          id: "piped-1",
          key: "piped-key",
          isNewKey: false,
        },
      });
    });
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

  describe("recreatePipedKey", () => {
    it(`should handle ${addNewPipedKey.fulfilled.type}`, () => {
      expect(
        pipedsSlice.reducer(baseState, {
          type: addNewPipedKey.fulfilled.type,
          payload: "add-new-piped-key",
          meta: {
            arg: {
              pipedId: "piped-1",
            },
          },
        })
      ).toEqual({
        ...baseState,
        registeredPiped: {
          id: "piped-1",
          key: "add-new-piped-key",
          isNewKey: true,
        },
      });
    });
  });

  describe("editPiped", () => {
    it(`should handle ${editPiped.pending.type}`, () => {
      expect(
        pipedsSlice.reducer(baseState, {
          type: editPiped.pending.type,
        })
      ).toEqual({
        ...baseState,
        updating: true,
      });
    });

    it(`should handle ${editPiped.rejected.type}`, () => {
      expect(
        pipedsSlice.reducer(
          { ...baseState, updating: true },
          {
            type: editPiped.rejected.type,
          }
        )
      ).toEqual({
        ...baseState,
        updating: false,
      });
    });

    it(`should handle ${editPiped.fulfilled.type}`, () => {
      expect(
        pipedsSlice.reducer(
          { ...baseState, updating: true },
          {
            type: editPiped.fulfilled.type,
          }
        )
      ).toEqual({
        ...baseState,
        updating: false,
      });
    });
  });
});
