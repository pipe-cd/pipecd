import {
  deploymentFrequencySlice,
  DeploymentFrequencyState,
  fetchDeploymentFrequency,
  fetchDeployment24h,
} from "./";

const initialState: DeploymentFrequencyState = {
  status: "idle",
  data: [],
  status24h: "idle",
  data24h: [],
};

describe("deploymentFrequencySlice reducer", () => {
  it("should return the initial state", () => {
    expect(
      deploymentFrequencySlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toEqual(initialState);
  });

  describe("fetchDeploymentFrequency", () => {
    it(`should handle ${fetchDeploymentFrequency.pending.type}`, () => {
      expect(
        deploymentFrequencySlice.reducer(initialState, {
          type: fetchDeploymentFrequency.pending.type,
        })
      ).toEqual({ ...initialState, status: "loading" });
    });

    it(`should handle ${fetchDeploymentFrequency.rejected.type}`, () => {
      expect(
        deploymentFrequencySlice.reducer(
          { ...initialState, status: "loading" },
          {
            type: fetchDeploymentFrequency.rejected.type,
          }
        )
      ).toEqual({ ...initialState, status: "failed" });
    });

    it(`should handle ${fetchDeploymentFrequency.fulfilled.type}`, () => {
      expect(
        deploymentFrequencySlice.reducer(
          { ...initialState, status: "loading" },
          {
            type: fetchDeploymentFrequency.fulfilled.type,
            payload: [],
          }
        )
      ).toEqual({ ...initialState, data: [], status: "succeeded" });
    });
  });

  describe("fetchDeploymentFrequency24h", () => {
    it(`should handle ${fetchDeployment24h.pending.type}`, () => {
      expect(
        deploymentFrequencySlice.reducer(initialState, {
          type: fetchDeployment24h.pending.type,
        })
      ).toEqual({ ...initialState, status24h: "loading" });
    });

    it(`should handle ${fetchDeployment24h.rejected.type}`, () => {
      expect(
        deploymentFrequencySlice.reducer(
          { ...initialState, status24h: "loading" },
          {
            type: fetchDeployment24h.rejected.type,
          }
        )
      ).toEqual({ ...initialState, status24h: "failed" });
    });

    it(`should handle ${fetchDeployment24h.fulfilled.type}`, () => {
      expect(
        deploymentFrequencySlice.reducer(
          { ...initialState, status24h: "loading" },
          {
            type: fetchDeployment24h.fulfilled.type,
            payload: [],
          }
        )
      ).toEqual({ ...initialState, data24h: [], status24h: "succeeded" });
    });
  });
});
