import { action } from "@storybook/addon-actions";
import * as React from "react";
import { GenerateAPIKeyDialog } from "./";

export default {
  title: "SETTINGS/APIKey/GenerateAPIKeyDialog",
  component: GenerateAPIKeyDialog,
};

export const overview: React.FC = () => (
  <GenerateAPIKeyDialog
    open
    onClose={action("onClose")}
    onSubmit={action("onSubmit")}
  />
);
