import { DeploymentStatus } from "./deployments";
import { ApplicationKind } from "./applications";
import {
  deploymentFilterOptionsSlice,
  updateDeploymentFilter,
  clearDeploymentFilter,
} from "./deployment-filter-options";

const initialState = {
  applicationIds: [],
  envIds: [],
  kinds: [],
  statuses: [],
};

describe("deploymentFilterOptionsSlice reducer", () => {
  it("should handle initial state", () => {
    expect(
      deploymentFilterOptionsSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toEqual(initialState);
  });

  it(`should handle ${updateDeploymentFilter.type}`, () => {
    expect(
      deploymentFilterOptionsSlice.reducer(initialState, {
        type: updateDeploymentFilter.type,
        payload: {
          applicationIds: ["app-1"],
          envIds: ["env-1"],
          kinds: [ApplicationKind.KUBERNETES],
          statuses: [DeploymentStatus.DEPLOYMENT_SUCCESS],
        },
      })
    ).toEqual({
      applicationIds: ["app-1"],
      envIds: ["env-1"],
      kinds: [ApplicationKind.KUBERNETES],
      statuses: [DeploymentStatus.DEPLOYMENT_SUCCESS],
    });
  });

  it(`should handle ${clearDeploymentFilter.type}`, () => {
    expect(
      deploymentFilterOptionsSlice.reducer(
        {
          applicationIds: ["app-1"],
          envIds: ["env-1"],
          kinds: [ApplicationKind.KUBERNETES],
          statuses: [DeploymentStatus.DEPLOYMENT_SUCCESS],
        },
        {
          type: clearDeploymentFilter.type,
        }
      )
    ).toEqual(initialState);
  });
});
