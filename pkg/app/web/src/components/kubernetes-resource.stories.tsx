import React from "react";
import { KubernetesResource } from "./kubernetes-resource";
import { HealthStatus } from "../modules/applications-live-state";
import { action } from "@storybook/addon-actions";

export default {
  title: "APPLICATION/KubernetesResource",
  component: KubernetesResource,
};

export const overview: React.FC = () => (
  <KubernetesResource
    resource={{
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
    }}
    onClick={action("onClick")}
  />
);
