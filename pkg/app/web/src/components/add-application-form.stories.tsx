import React from "react";
import { AddApplicationForm } from "./add-application-form";
import { action } from "@storybook/addon-actions";
import { createDecoratorRedux } from "../../.storybook/redux-decorator";

export default {
  title: "AddApplicationForm",
  component: AddApplicationForm,
  decorators: [createDecoratorRedux({})],
};

export const overview: React.FC = () => (
  <AddApplicationForm
    projectName="pipe-cd"
    onSubmit={action("onSubmit")}
    onClose={action("onClose")}
    isAdding={false}
  />
);
