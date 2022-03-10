import { dummyApplicationLiveState } from "~/__fixtures__/dummy-application-live-state";
import {
  applicationLiveStateSlice,
  ApplicationLiveStateState,
  fetchApplicationStateById,
} from "./";

const initialState: ApplicationLiveStateState = {
  entities: {},
  hasError: {},
  loading: {},
  ids: [],
};

describe("applicationLiveStateSlice reducer", () => {
  it("should return the initial state", () => {
    expect(
      applicationLiveStateSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toEqual(initialState);
  });

  describe("fetchApplicationStateById", () => {
    it(`should handle ${fetchApplicationStateById.pending.type}`, () => {
      expect(
        applicationLiveStateSlice.reducer(initialState, {
          type: fetchApplicationStateById.pending.type,
          meta: {
            arg: "application-1",
          },
        })
      ).toEqual({
        ...initialState,
        hasError: { "application-1": false },
        loading: { "application-1": true },
      });
    });

    it(`should handle ${fetchApplicationStateById.rejected.type}`, () => {
      expect(
        applicationLiveStateSlice.reducer(
          { ...initialState, hasError: { "application-1": false } },
          {
            type: fetchApplicationStateById.rejected.type,
            meta: {
              arg: "application-1",
            },
          }
        )
      ).toEqual({
        ...initialState,
        hasError: { "application-1": true },
        loading: { "application-1": false },
      });
    });

    it(`should handle ${fetchApplicationStateById.fulfilled.type}`, () => {
      expect(
        applicationLiveStateSlice.reducer(
          { ...initialState, hasError: { "application-1": false } },
          {
            type: fetchApplicationStateById.fulfilled.type,
            meta: {
              arg: "application-1",
            },
            payload: dummyApplicationLiveState,
          }
        )
      ).toEqual({
        entities: {
          [dummyApplicationLiveState.applicationId]: dummyApplicationLiveState,
        },
        ids: [dummyApplicationLiveState.applicationId],
        hasError: { "application-1": false },
        loading: { "application-1": false },
      });
    });
  });
});
