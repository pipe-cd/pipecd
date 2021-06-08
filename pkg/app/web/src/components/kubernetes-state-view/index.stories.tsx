import { Story } from "@storybook/react";
import { resourcesList } from "../../__fixtures__/dummy-application-live-state";
import { KubernetesStateView, KubernetesStateViewProps } from "./";

export default {
  title: "APPLICATION/KubernetesStateView",
  component: KubernetesStateView,
};

const Template: Story<KubernetesStateViewProps> = (args) => (
  <KubernetesStateView {...args} />
);
export const Overview = Template.bind({});
Overview.args = { resources: resourcesList };
