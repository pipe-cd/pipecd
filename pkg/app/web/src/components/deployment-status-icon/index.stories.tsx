import { DeploymentStatusIcon, DeploymentStatusIconProps } from "./";
import { Story } from "@storybook/react";
import { DeploymentStatus } from "~/modules/deployments";

export default {
  title: "DEPLOYMENT/StatusIcon",
  component: DeploymentStatusIcon,
};

const Template: Story<DeploymentStatusIconProps> = (args) => (
  <DeploymentStatusIcon {...args} />
);

export const Success = Template.bind({});
Success.args = { status: DeploymentStatus.DEPLOYMENT_SUCCESS };

export const Failure = Template.bind({});
Failure.args = { status: DeploymentStatus.DEPLOYMENT_FAILURE };

export const Cancelled = Template.bind({});
Cancelled.args = { status: DeploymentStatus.DEPLOYMENT_CANCELLED };

export const Pending = Template.bind({});
Pending.args = { status: DeploymentStatus.DEPLOYMENT_PENDING };

export const Planned = Template.bind({});
Planned.args = { status: DeploymentStatus.DEPLOYMENT_PLANNED };

export const Running = Template.bind({});
Running.args = { status: DeploymentStatus.DEPLOYMENT_RUNNING };
