import React from "react";
import { DeploymentItem } from "./deployment-item";

export default {
  title: "DeploymentItem",
  component: DeploymentItem,
};

export const overview: React.FC = () => <DeploymentItem id="app-1" />;
