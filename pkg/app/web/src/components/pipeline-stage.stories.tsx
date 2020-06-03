import React from "react";
import { PipelineStage } from "./pipeline-stage";
import { StageStatus } from "pipe/pkg/app/web/model/deployment_pb";

export default {
  title: "PipelineStage",
  component: PipelineStage
};

export const overview: React.FC = () => (
  <PipelineStage status={StageStatus.STAGE_SUCCESS} name="K8S_CANARY_ROLLOUT" />
);
