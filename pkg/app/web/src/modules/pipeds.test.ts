import { dummyPiped } from "../__fixtures__/dummy-piped";
import {
  pipedsSlice,
  clearRegisteredPipedInfo,
  addPiped,
  fetchPipeds,
  recreatePipedKey,
  selectPipedsByEnv,
  Piped,
} from "./pipeds";

test("selectPipedsByEnv", () => {
  const disabledPiped: Piped = { ...dummyPiped, id: "piped-2", disabled: true };
  expect(selectPipedsByEnv({ entities: {}, ids: [] }, "env-1")).toEqual([]);
  expect(
    selectPipedsByEnv(
      {
        entities: {
          [dummyPiped.id]: dummyPiped,
          [disabledPiped.id]: disabledPiped,
        },
        ids: [dummyPiped.id, disabledPiped.id],
      },
      dummyPiped.envIdsList[0]
    )
  ).toEqual([dummyPiped]);
});

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
