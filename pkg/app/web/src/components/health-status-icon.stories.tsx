import React from "react";
import { KubernetesResourceHealthStatusIcon } from "./health-status-icon";
import { HealthStatus } from "../modules/applications-live-state";

export default {
  title: "APPLICATION|HealthStatusIcon",
  component: KubernetesResourceHealthStatusIcon,
};

export const overview: React.FC = () => (
  <KubernetesResourceHealthStatusIcon health={HealthStatus.HEALTHY} />
);
