import React from "react";
import { KubernetesResource } from "./kubernetes-resource";
import { HealthStatus } from "../modules/applications-live-state";

export default {
  title: "APPLICATION|KubernetesResource",
  component: KubernetesResource,
};

export const overview: React.FC = () => (
  <KubernetesResource
    name="demo-application"
    kind="Ingress"
    health={HealthStatus.HEALTHY}
  />
);
