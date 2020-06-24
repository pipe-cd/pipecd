import React from "react";
import { ApplicationList } from "./application-list";

export default {
  title: "ApplicationList",
  component: ApplicationList,
};

export const overview: React.FC = () => <ApplicationList />;
