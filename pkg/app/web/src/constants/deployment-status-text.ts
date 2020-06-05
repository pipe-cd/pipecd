import { DeploymentStatus } from "pipe/pkg/app/web/model/deployment_pb";

export const DEPLOYMENT_STATE_TEXT: Record<DeploymentStatus, string> = {
  [DeploymentStatus.DEPLOYMENT_PENDING]: "Pending",
  [DeploymentStatus.DEPLOYMENT_PLANNED]: "Planned",
  [DeploymentStatus.DEPLOYMENT_RUNNING]: "Running",
  [DeploymentStatus.DEPLOYMENT_ROLLING_BACK]: "Rolling Back",
  [DeploymentStatus.DEPLOYMENT_SUCCESS]: "Success",
  [DeploymentStatus.DEPLOYMENT_FAILURE]: "Failure",
  [DeploymentStatus.DEPLOYMENT_CANCELLED]: "Canceled"
};
