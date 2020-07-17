import React from "react";
import { AddEnvForm } from "./add-env-form";
import { action } from "@storybook/addon-actions";

export default {
  title: "SETTINGS|AddEnvForm",
  component: AddEnvForm,
};

export const overview: React.FC = () => (
  <AddEnvForm
    onClose={action("onClose")}
    onSubmit={action("onSubmit")}
    projectName="project-name"
  />
);
