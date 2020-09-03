import React from "react";
import { InputForm } from "./input-form";
import { action } from "@storybook/addon-actions";

export default {
  title: "InputForm",
  component: InputForm,
};

export const overview: React.FC = () => (
  <InputForm name="Name" currentValue="value" onSave={action("onChange")} />
);

export const IsSecret: React.FC = () => (
  <InputForm
    name="Name"
    currentValue="value"
    onSave={action("onChange")}
    isSecret
  />
);
