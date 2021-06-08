import { Story } from "@storybook/react";
import { KubernetesResourceDetail, KubernetesResourceDetailProps } from "./";

export default {
  title: "KubernetesResourceDetail",
  component: KubernetesResourceDetail,
  argTypes: {
    onClose: {
      action: "onClose",
    },
  },
};

const Template: Story<KubernetesResourceDetailProps> = (args) => (
  <KubernetesResourceDetail {...args} />
);
export const Overview = Template.bind({});
Overview.args = {
  resource: {
    name: "demo-application-9504e8601a",
    namespace: "default",
    apiVersion: "apps/v1",
    healthDescription: "Unimplemented or unknown resource",
    kind: "Pod",
  },
};
