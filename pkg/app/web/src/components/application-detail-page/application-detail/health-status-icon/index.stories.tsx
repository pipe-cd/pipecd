import { Story } from "@storybook/react";
import { ApplicationLiveStateSnapshot } from "~/modules/applications-live-state";
import {
  ApplicationHealthStatusIcon,
  ApplicationHealthStatusIconProps,
} from ".";

export default {
  title: "APPLICATION/ApplicationHealthStatusIcon",
  component: ApplicationHealthStatusIcon,
};

const Template: Story<ApplicationHealthStatusIconProps> = (args) => (
  <ApplicationHealthStatusIcon {...args} />
);
export const Overview = Template.bind({});
Overview.args = { health: ApplicationLiveStateSnapshot.Status.HEALTHY };
