import { KubernetesResource, KubernetesResourceProps } from "./";
import { HealthStatus } from "../../modules/applications-live-state";
import { Story } from "@storybook/react";

export default {
  title: "APPLICATION/KubernetesResource",
  component: KubernetesResource,
  argTypes: {
    onClick: {
      action: "onClick",
    },
  },
};

const Template: Story<KubernetesResourceProps> = (args) => (
  <KubernetesResource {...args} />
);
export const Overview = Template.bind({});
Overview.args = {
  resource: {
    apiVersion: "v1",
    createdAt: 0,
    healthStatus: HealthStatus.HEALTHY,
    healthDescription: "",
    id: "1",
    kind: "Pod",
    name: "resource-1",
    namespace: "default",
    ownerIdsList: [],
    parentIdsList: [],
    updatedAt: 0,
  },
};
