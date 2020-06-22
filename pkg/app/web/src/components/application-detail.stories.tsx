import React from "react";
import { ApplicationDetail } from "./application-detail";

export default {
  title: "ApplicationDetail",
  component: ApplicationDetail,
};

export const overview: React.FC = () => (
  <ApplicationDetail applicationId="application-1" />
);
