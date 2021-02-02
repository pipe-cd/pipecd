import React from "react";
import { AddEnvForm } from "./";
import { action } from "@storybook/addon-actions";

export default {
  title: "SETTINGS/AddEnvForm",
  component: AddEnvForm,
};

export const overview: React.FC = () => (
  <AddEnvForm
    onCancel={action("onCancel")}
    onSubmit={action("onSubmit")}
    projectName="project-name"
  />
);
