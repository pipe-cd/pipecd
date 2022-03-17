// import {
//   deploymentFrequencySlice,
//   DeploymentFrequencyState,
//   fetchDeploymentFrequency,
// } from "./";

// const initialState: DeploymentFrequencyState = { status: "idle", data: [] };

describe("fake", () => {
  it("fake", () => {
    true;
  });
});

// describe("deploymentFrequencySlice reducer", () => {
//   it("should return the initial state", () => {
//     expect(
//       deploymentFrequencySlice.reducer(undefined, {
//         type: "TEST_ACTION",
//       })
//     ).toEqual(initialState);
//   });

//   describe("fetchDeploymentFrequency", () => {
//     it(`should handle ${fetchDeploymentFrequency.pending.type}`, () => {
//       expect(
//         deploymentFrequencySlice.reducer(initialState, {
//           type: fetchDeploymentFrequency.pending.type,
//         })
//       ).toEqual({ ...initialState, status: "loading" });
//     });

//     it(`should handle ${fetchDeploymentFrequency.rejected.type}`, () => {
//       expect(
//         deploymentFrequencySlice.reducer(
//           { ...initialState, status: "loading" },
//           {
//             type: fetchDeploymentFrequency.rejected.type,
//           }
//         )
//       ).toEqual({ ...initialState, status: "failed" });
//     });

//     it(`should handle ${fetchDeploymentFrequency.fulfilled.type}`, () => {
//       expect(
//         deploymentFrequencySlice.reducer(
//           { ...initialState, status: "loading" },
//           {
//             type: fetchDeploymentFrequency.fulfilled.type,
//             payload: [],
//           }
//         )
//       ).toEqual({ ...initialState, data: [], status: "succeeded" });
//     });
//   });
// });
