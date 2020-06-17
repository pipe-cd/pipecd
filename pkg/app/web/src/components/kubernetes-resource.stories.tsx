import React from "react";
import { KubernetesResource } from "./kubernetes-resource";

export default {
  title: "KubernetesResource",
  component: KubernetesResource,
};

export const overview: React.FC = () => (
  <KubernetesResource name="demo-application" kind="Ingress" />
);
