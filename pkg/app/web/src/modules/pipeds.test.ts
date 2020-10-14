import { dummyPiped } from "../__fixtures__/dummy-piped";
import {
  pipedsSlice,
  clearRegisteredPipedInfo,
  addPiped,
  fetchPipeds,
  recreatePipedKey,
} from "./pipeds";

describe("pipedsSlice reducer", () => {
  it("should handle initial state", () => {
    expect(
      pipedsSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toEqual({
      entities: {},
      ids: [],
      registeredPiped: null,
    });
  });

  it(`should handle ${clearRegisteredPipedInfo.type}`, () => {
    expect(
      pipedsSlice.reducer(
        {
          entities: {},
          ids: [],
          registeredPiped: {
            id: "piped-1",
            key: "piped-key",
          },
        },
        {
          type: clearRegisteredPipedInfo.type,
        }
      )
    ).toEqual({
      entities: {},
      ids: [],
      registeredPiped: null,
    });
  });

  describe("addPiped", () => {
    it(`should handle ${addPiped.fulfilled.type}`, () => {
      expect(
        pipedsSlice.reducer(
          {
            entities: {},
            ids: [],
            registeredPiped: null,
          },
          {
            type: addPiped.fulfilled.type,
            payload: {
              id: "piped-1",
              key: "piped-key",
            },
          }
        )
      ).toEqual({
        entities: {},
        ids: [],
        registeredPiped: {
          id: "piped-1",
          key: "piped-key",
        },
      });
    });
  });

  describe("fetchPipeds", () => {
    it(`should handle ${fetchPipeds.fulfilled.type}`, () => {
      expect(
        pipedsSlice.reducer(
          {
            entities: {},
            ids: [],
            registeredPiped: null,
          },
          {
            type: fetchPipeds.fulfilled.type,
            payload: [dummyPiped],
          }
        )
      ).toEqual({
        entities: { [dummyPiped.id]: dummyPiped },
        ids: [dummyPiped.id],
        registeredPiped: null,
      });
    });
  });

  describe("recreatePipedKey", () => {
    it(`should handle ${recreatePipedKey.fulfilled.type}`, () => {
      expect(
        pipedsSlice.reducer(
          {
            entities: {},
            ids: [],
            registeredPiped: null,
          },
          {
            type: recreatePipedKey.fulfilled.type,
            payload: "recreated-piped-key",
            meta: {
              arg: {
                pipedId: "piped-1",
              },
            },
          }
        )
      ).toEqual({
        entities: {},
        ids: [],
        registeredPiped: {
          id: "piped-1",
          key: "recreated-piped-key",
        },
      });
    });
  });
});
