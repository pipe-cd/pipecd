import React from "react";
import { DeploymentConfigForm } from "./deployment-config-form";
import { createDecoratorRedux } from "../../.storybook/redux-decorator";
import { action } from "@storybook/addon-actions";

export default {
  title: "DeploymentConfigForm",
  component: DeploymentConfigForm,
  decorators: [createDecoratorRedux({})],
};

export const overview: React.FC = () => (
  <DeploymentConfigForm
    applicationId="application-1"
    onSkip={action("onSkip")}
  />
);
