import React from "react";
import { DeploymentConfigForm } from "./";
import { action } from "@storybook/addon-actions";
import { createStore } from "../../../test-utils";
import { Provider } from "react-redux";
import { dummyDeploymentConfigTemplates } from "../../__fixtures__/dummy-deployment-config";

export default {
  title: "DEPLOYMENT/DeploymentConfigForm",
  component: DeploymentConfigForm,
};

export const overview: React.FC = () => (
  <Provider
    store={createStore({
      deploymentConfigs: {
        targetApplicationId: "application-1",
        templates: {
          "application-1": dummyDeploymentConfigTemplates,
        },
      },
    })}
  >
    <DeploymentConfigForm
      applicationId="application-1"
      onSkip={action("onSkip")}
    />
  </Provider>
);

export const loading: React.FC = () => (
  <Provider store={createStore({})}>
    <DeploymentConfigForm
      applicationId="application-1"
      onSkip={action("onSkip")}
    />
  </Provider>
);
