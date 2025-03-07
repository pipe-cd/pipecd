import {
  deploymentChangeFailureRateSlice,
  DeploymentChangeFailureRateState,
  fetchDeploymentChangeFailureRate,
  fetchDeploymentChangeFailureRate24h,
} from "./";

const initialState: DeploymentChangeFailureRateState = {
  status: "idle",
  data: [],
  status24h: "idle",
  data24h: [],
};

describe("deploymentChangeFailureRateSlice reducer", () => {
  it("should return the initial state", () => {
    expect(
      deploymentChangeFailureRateSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toEqual(initialState);
  });

  describe("fetchDeploymentChangeFailureRate", () => {
    it(`should handle ${fetchDeploymentChangeFailureRate.pending.type}`, () => {
      expect(
        deploymentChangeFailureRateSlice.reducer(initialState, {
          type: fetchDeploymentChangeFailureRate.pending.type,
        })
      ).toEqual({ ...initialState, status: "loading" });
    });

    it(`should handle ${fetchDeploymentChangeFailureRate.rejected.type}`, () => {
      expect(
        deploymentChangeFailureRateSlice.reducer(
          { ...initialState, status: "loading" },
          {
            type: fetchDeploymentChangeFailureRate.rejected.type,
          }
        )
      ).toEqual({ ...initialState, status: "failed" });
    });

    it(`should handle ${fetchDeploymentChangeFailureRate.fulfilled.type}`, () => {
      expect(
        deploymentChangeFailureRateSlice.reducer(
          { ...initialState, status: "loading" },
          {
            type: fetchDeploymentChangeFailureRate.fulfilled.type,
            payload: [],
          }
        )
      ).toEqual({ ...initialState, data: [], status: "succeeded" });
    });
  });

  describe("fetchDeploymentChangeFailureRate24h", () => {
    it(`should handle ${fetchDeploymentChangeFailureRate24h.pending.type}`, () => {
      expect(
        deploymentChangeFailureRateSlice.reducer(initialState, {
          type: fetchDeploymentChangeFailureRate24h.pending.type,
        })
      ).toEqual({ ...initialState, status24h: "loading" });
    });

    it(`should handle ${fetchDeploymentChangeFailureRate24h.rejected.type}`, () => {
      expect(
        deploymentChangeFailureRateSlice.reducer(
          { ...initialState, status24h: "loading" },
          {
            type: fetchDeploymentChangeFailureRate24h.rejected.type,
          }
        )
      ).toEqual({ ...initialState, status24h: "failed" });
    });

    it(`should handle ${fetchDeploymentChangeFailureRate24h.fulfilled.type}`, () => {
      expect(
        deploymentChangeFailureRateSlice.reducer(
          { ...initialState, status24h: "loading" },
          {
            type: fetchDeploymentChangeFailureRate24h.fulfilled.type,
            payload: [],
          }
        )
      ).toEqual({ ...initialState, data24h: [], status24h: "succeeded" });
    });
  });
});
