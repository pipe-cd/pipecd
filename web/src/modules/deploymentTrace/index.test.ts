import { LoadingStatus } from "~/types/module";
import {
  deploymentTraceSlice,
  fetchDeploymentTraces,
  fetchMoreDeploymentTraces,
} from "./index";
import { dummyDeploymentTrace } from "~/__fixtures__/dummy-deployment-trace";

describe("deploymentTrace slice", () => {
  const initialState = {
    status: "idle" as LoadingStatus,
    loading: {},
    hasMore: true,
    cursor: "",
    minUpdatedAt: 0,
    skippable: {},
    canceling: {},
    entities: {},
    ids: [],
  };

  it("should return the initial state", () => {
    expect(
      deploymentTraceSlice.reducer(initialState, {
        type: "TEST_ACTION",
      })
    ).toEqual(initialState);
  });

  describe("fetchDeploymentTrace", () => {
    it(`should handle ${fetchDeploymentTraces.pending.type}`, () => {
      expect(
        deploymentTraceSlice.reducer(initialState, {
          type: fetchDeploymentTraces.pending.type,
        })
      ).toEqual({
        ...initialState,
        status: "loading",
      });
    });

    it(`should handle ${fetchDeploymentTraces.rejected.type}`, () => {
      expect(
        deploymentTraceSlice.reducer(
          {
            ...initialState,
            status: "loading",
          },
          {
            type: fetchDeploymentTraces.rejected.type,
          }
        )
      ).toEqual({
        ...initialState,
        status: "failed",
      });
    });

    it(`should handle ${fetchDeploymentTraces.fulfilled.type}`, () => {
      expect(
        deploymentTraceSlice.reducer(
          {
            ...initialState,
            status: "loading",
          },
          {
            type: fetchDeploymentTraces.fulfilled.type,
            payload: {
              tracesList: [dummyDeploymentTrace],
              cursor: "next cursor",
            },
          }
        )
      ).toEqual({
        ...initialState,
        entities: dummyDeploymentTrace.trace
          ? { [dummyDeploymentTrace.trace.id]: dummyDeploymentTrace }
          : {},
        hasMore: false,
        ids: dummyDeploymentTrace.trace?.id
          ? [dummyDeploymentTrace.trace.id]
          : [],
        status: "succeeded",
        cursor: "next cursor",
        minUpdatedAt: 0,
      });
    });
  });

  describe("fetchMoreDeploymentTraces", () => {
    it(`should handle ${fetchMoreDeploymentTraces.pending.type}`, () => {
      expect(
        deploymentTraceSlice.reducer(initialState, {
          type: fetchMoreDeploymentTraces.pending.type,
        })
      ).toEqual({ ...initialState, status: "loading" });
    });

    it(`should handle ${fetchMoreDeploymentTraces.rejected.type}`, () => {
      expect(
        deploymentTraceSlice.reducer(
          { ...initialState, status: "loading" },
          {
            type: fetchMoreDeploymentTraces.rejected.type,
          }
        )
      ).toEqual({ ...initialState, status: "failed" });
    });

    it(`should handle ${fetchMoreDeploymentTraces.fulfilled.type}`, () => {
      expect(
        deploymentTraceSlice.reducer(
          { ...initialState, status: "loading" },
          {
            type: fetchMoreDeploymentTraces.fulfilled.type,
            payload: {
              tracesList: [dummyDeploymentTrace],
              cursor: "next cursor",
            },
          }
        )
      ).toEqual({
        ...initialState,
        hasMore: false,
        ids: dummyDeploymentTrace.trace?.id
          ? [dummyDeploymentTrace.trace.id]
          : [],
        entities: dummyDeploymentTrace.trace
          ? { [dummyDeploymentTrace.trace.id]: dummyDeploymentTrace }
          : {},
        status: "succeeded",
        cursor: "next cursor",
        minUpdatedAt: -2592000,
      });
    });
  });
});
