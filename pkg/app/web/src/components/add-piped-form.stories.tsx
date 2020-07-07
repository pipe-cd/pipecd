import React from "react";
import { AddPipedForm } from "./add-piped-form";
import { action } from "@storybook/addon-actions";

export default {
  title: "SETTINGS|AddPipedForm",
  component: AddPipedForm,
};

export const overview: React.FC = () => (
  <AddPipedForm
    onClose={action("onClose")}
    onSubmit={action("onSubmit")}
    projectName="project-name"
  />
);
