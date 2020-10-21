import React from "react";
import { resourcesList } from "../__fixtures__/dummy-application-live-state";
import { KubernetesStateView } from "./kubernetes-state-view";

export default {
  title: "APPLICATION/KubernetesStateView",
  component: KubernetesStateView,
};

export const overview: React.FC = () => (
  <KubernetesStateView
    showKinds={["ReplicaSet", "Pod"]}
    resources={resourcesList}
  />
);
