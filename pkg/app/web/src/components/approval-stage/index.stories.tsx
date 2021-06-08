import { Story } from "@storybook/react";
import { ApprovalStage, ApprovalStageProps } from "./";

export default {
  title: "DEPLOYMENT/Pipeline/ApprovalStage",
  component: ApprovalStage,
  argTypes: {
    onClick: { action: "onClick" },
  },
};

const Template: Story<ApprovalStageProps> = (args) => (
  <ApprovalStage {...args} />
);
export const Overview = Template.bind({});
Overview.args = { id: "stage-1", name: "K8S_CANARY_ROLLOUT", active: false };
