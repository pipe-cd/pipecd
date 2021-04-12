import * as React from "react";
import { resourcesList } from "../../__fixtures__/dummy-application-live-state";
import { KubernetesStateView } from "./";

export default {
  title: "APPLICATION/KubernetesStateView",
  component: KubernetesStateView,
};

export const overview: React.FC = () => (
  <KubernetesStateView resources={resourcesList} />
);
