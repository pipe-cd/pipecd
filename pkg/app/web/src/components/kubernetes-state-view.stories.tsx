import React from "react";
import { KubernetesStateView } from "./kubernetes-state-view";

export default {
  title: "KubernetesStateView",
  component: KubernetesStateView
};

export const overview: React.FC = () => (
  <KubernetesStateView />
);