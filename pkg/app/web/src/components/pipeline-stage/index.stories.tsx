import { Story } from "@storybook/react";
import { StageStatus } from "pipe/pkg/app/web/model/deployment_pb";
import { PipelineStage, PipelineStageProps } from "./";

export default {
  title: "DEPLOYMENT/Pipeline/PipelineStage",
  component: PipelineStage,
  argTypes: {
    onClick: {
      action: "onClick",
    },
  },
};

const Template: Story<PipelineStageProps> = (args) => (
  <PipelineStage {...args} />
);
export const Overview = Template.bind({});
Overview.args = {
  id: "stage-1",
  status: StageStatus.STAGE_SUCCESS,
  name: "K8S_CANARY_ROLLOUT",
  active: false,
  metadata: [],
  isDeploymentRunning: true,
};

export const LongName = Template.bind({});
LongName.args = {
  id: "stage-1",
  status: StageStatus.STAGE_SUCCESS,
  name: "LONG_STAGE_NAME_XXXXXXX_YYYYYY_ZZZZZZZZ",
  active: false,
  metadata: [],
  isDeploymentRunning: true,
};

export const Stopped = Template.bind({});
Stopped.args = {
  id: "stage-1",
  status: StageStatus.STAGE_NOT_STARTED_YET,
  name: "K8S_CANARY_ROLLOUT",
  active: false,
  metadata: [],
  isDeploymentRunning: false,
};

export const Approved = Template.bind({});
Approved.args = {
  id: "stage-1",
  status: StageStatus.STAGE_SUCCESS,
  name: "K8S_CANARY_ROLLOUT",
  active: false,
  metadata: [],
  approver: "User",
  isDeploymentRunning: true,
};

export const TrafficPercentage = Template.bind({});
TrafficPercentage.args = {
  id: "stage-1",
  status: StageStatus.STAGE_SUCCESS,
  name: "K8S_CANARY_ROLLOUT",
  active: false,
  metadata: [
    ["baseline-percentage", "0"],
    ["canary-percentage", "50"],
    ["primary-percentage", "50"],
  ],
  isDeploymentRunning: true,
};

export const PromotePercentage = Template.bind({});
PromotePercentage.args = {
  id: "stage-1",
  status: StageStatus.STAGE_SUCCESS,
  name: "K8S_CANARY_ROLLOUT",
  active: false,
  metadata: [["promote-percentage", "75"]],
  isDeploymentRunning: true,
};
