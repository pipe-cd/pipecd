import { Story } from "@storybook/react";
import { StageStatus } from "~/modules/deployments";
import { StageStatusIcon, StageStatusIconProps } from ".";

export default {
  title: "DEPLOYMENT/StageStatusIcon",
  component: StageStatusIcon,
};

const Template: Story<StageStatusIconProps> = (args) => (
  <StageStatusIcon {...args} />
);

export const Cancelled = Template.bind({});
Cancelled.args = { status: StageStatus.STAGE_CANCELLED };

export const Failure = Template.bind({});
Failure.args = { status: StageStatus.STAGE_FAILURE };

export const NotStartedYet = Template.bind({});
NotStartedYet.args = { status: StageStatus.STAGE_NOT_STARTED_YET };

export const Running = Template.bind({});
Running.args = { status: StageStatus.STAGE_RUNNING };

export const Success = Template.bind({});
Success.args = { status: StageStatus.STAGE_SUCCESS };
