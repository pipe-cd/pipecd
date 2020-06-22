import React from "react";
import { ApplicationStateView } from "./application-state-view";

export default {
  title: "ApplicationStateView",
  component: ApplicationStateView,
};

export const overview: React.FC = () => (
  <ApplicationStateView applicationId="application-1" />
);
