import {
  deploymentChangeFailureRateSlice,
  DeploymentChangeFailureRateState,
  fetchDeploymentChangeFailureRate,
} from "./";

const initialState: DeploymentChangeFailureRateState = {
  status: "idle",
  data: [],
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
});
