import { Story } from "@storybook/react";
import { HealthStatus } from "~/modules/applications-live-state";
import {
  KubernetesResourceHealthStatusIcon,
  KubernetesResourceHealthStatusIconProps,
} from ".";

export default {
  title: "APPLICATION/KubernetesResourceHealthStatusIcon",
  component: KubernetesResourceHealthStatusIcon,
};

const Template: Story<KubernetesResourceHealthStatusIconProps> = (args) => (
  <KubernetesResourceHealthStatusIcon {...args} />
);
export const Overview = Template.bind({});
Overview.args = { health: HealthStatus.HEALTHY };
