import React from "react";
import { PipelineStage } from "./pipeline-stage";
import { StageStatus } from "pipe/pkg/app/web/model/deployment_pb";
import { action } from "@storybook/addon-actions";

export default {
  title: "DEPLOYMENT|PipelineStage",
  component: PipelineStage,
};

export const overview: React.FC = () => (
  <PipelineStage
    id="stage-1"
    status={StageStatus.STAGE_SUCCESS}
    name="K8S_CANARY_ROLLOUT"
    onClick={action("onClick")}
    active={false}
  />
);
