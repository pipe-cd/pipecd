import {
  DeploymentStatus,
  PipelineStage,
  Deployment,
  StageStatus,
} from "pipecd/web/model/deployment_pb";

export type Stage = Required<PipelineStage.AsObject>;
export type DeploymentStatusKey = keyof typeof DeploymentStatus;

export { Deployment, DeploymentStatus, StageStatus, PipelineStage };
