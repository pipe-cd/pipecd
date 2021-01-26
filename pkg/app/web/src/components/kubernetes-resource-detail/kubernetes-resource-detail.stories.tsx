import { action } from "@storybook/addon-actions";
import React from "react";
import { KubernetesResourceDetail } from "./kubernetes-resource-detail";

export default {
  title: "KubernetesResourceDetail",
  component: KubernetesResourceDetail,
};

export const overview: React.FC = () => (
  <KubernetesResourceDetail
    resource={{
      name: "demo-application-9504e8601a",
      namespace: "default",
      apiVersion: "apps/v1",
      healthDescription: "Unimplemented or unknown resource",
      kind: "Pod",
    }}
    onClose={action("onClose")}
  />
);
