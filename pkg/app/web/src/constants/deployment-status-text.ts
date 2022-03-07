import { DeploymentStatus } from "pipecd/pkg/app/web/model/deployment_pb";

export const DEPLOYMENT_STATE_TEXT: Record<DeploymentStatus, string> = {
  [DeploymentStatus.DEPLOYMENT_PENDING]: "PENDING",
  [DeploymentStatus.DEPLOYMENT_PLANNED]: "PLANNED",
  [DeploymentStatus.DEPLOYMENT_RUNNING]: "RUNNING",
  [DeploymentStatus.DEPLOYMENT_ROLLING_BACK]: "ROLLING BACK",
  [DeploymentStatus.DEPLOYMENT_SUCCESS]: "SUCCESS",
  [DeploymentStatus.DEPLOYMENT_FAILURE]: "FAILURE",
  [DeploymentStatus.DEPLOYMENT_CANCELLED]: "CANCELLED",
};
