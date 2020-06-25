import React from "react";
import { HealthStatusIcon } from "./health-status-icon";
import { HealthStatus } from "../modules/applications-live-state";

export default {
  title: "HealthStatusIcon",
  component: HealthStatusIcon,
};

export const overview: React.FC = () => (
  <HealthStatusIcon health={HealthStatus.HEALTHY} />
);
