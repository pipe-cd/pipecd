import {
  KubernetesResourceHealthStatusIcon,
  KubernetesResourceHealthStatusIconProps,
} from "./";
import { HealthStatus } from "../../modules/applications-live-state";
import { Story } from "@storybook/react";

export default {
  title: "APPLICATION/HealthStatusIcon",
  component: KubernetesResourceHealthStatusIcon,
};

const Template: Story<KubernetesResourceHealthStatusIconProps> = (args) => (
  <KubernetesResourceHealthStatusIcon {...args} />
);
export const Overview = Template.bind({});
Overview.args = { health: HealthStatus.HEALTHY };
