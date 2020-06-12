import React from "react";
import { ApplicationResource } from "./application-resource";

export default {
  title: "ApplicationResource",
  component: ApplicationResource,
};

export const overview: React.FC = () => (
  <ApplicationResource name="demo-application" kind="Ingress" />
);
