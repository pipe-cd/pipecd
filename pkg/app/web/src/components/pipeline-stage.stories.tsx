import React from "react";
import { PipelineStage } from "./pipeline-stage";
import { StageStatus } from "pipe/pkg/app/web/model/deployment_pb";
import { action } from "@storybook/addon-actions";

export default {
  title: "DEPLOYMENT/Pipeline/PipelineStage",
  component: PipelineStage,
};

export const overview: React.FC = () => (
  <PipelineStage
    id="stage-1"
    status={StageStatus.STAGE_SUCCESS}
    name="K8S_CANARY_ROLLOUT"
    onClick={action("onClick")}
    active={false}
    isDeploymentRunning
  />
);

export const longName: React.FC = () => (
  <PipelineStage
    id="stage-1"
    status={StageStatus.STAGE_SUCCESS}
    name="LONG_STAGE_NAME_XXXXXXX_YYYYYY_ZZZZZZZZ"
    onClick={action("onClick")}
    active={false}
    isDeploymentRunning
  />
);

export const stopped: React.FC = () => (
  <PipelineStage
    id="stage-1"
    status={StageStatus.STAGE_NOT_STARTED_YET}
    name="K8S_CANARY_ROLLOUT"
    onClick={action("onClick")}
    active={false}
    isDeploymentRunning={false}
  />
);

export const Approved: React.FC = () => (
  <PipelineStage
    id="stage-1"
    status={StageStatus.STAGE_SUCCESS}
    name="K8S_CANARY_ROLLOUT"
    onClick={action("onClick")}
    active={false}
    approver="User"
    isDeploymentRunning
  />
);
