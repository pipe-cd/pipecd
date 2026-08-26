import { DeploymentStatus } from "~~/model/deployment_pb";
import { isDeploymentRunning } from "./is-deployment-running";

describe("isDeploymentRunning", () => {
  it.each([
    [DeploymentStatus.DEPLOYMENT_PENDING, true],
    [DeploymentStatus.DEPLOYMENT_PLANNED, true],
    [DeploymentStatus.DEPLOYMENT_RUNNING, true],
    [DeploymentStatus.DEPLOYMENT_ROLLING_BACK, true],
    [DeploymentStatus.DEPLOYMENT_SUCCESS, false],
    [DeploymentStatus.DEPLOYMENT_FAILURE, false],
    [DeploymentStatus.DEPLOYMENT_CANCELLED, false],
    [undefined, false],
  ])("given status %p, returns %p", (status, expected) => {
    expect(isDeploymentRunning(status)).toBe(expected);
  });
});
