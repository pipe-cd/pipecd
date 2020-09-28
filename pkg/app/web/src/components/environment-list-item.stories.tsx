import React from "react";
import { EnvironmentListItem } from "./environment-list-item";

export default {
  title: "EnvironmentListItem",
  component: EnvironmentListItem,
};

export const overview: React.FC = () => <EnvironmentListItem id="env-id" />;
