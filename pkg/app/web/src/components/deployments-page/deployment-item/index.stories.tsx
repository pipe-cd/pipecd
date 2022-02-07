import { Story } from "@storybook/react";
import { Provider } from "react-redux";
import { dummyApplication } from "~/__fixtures__/dummy-application";
import { dummyDeployment } from "~/__fixtures__/dummy-deployment";
import { createStore } from "~~/test-utils";
import { DeploymentItem } from ".";

export default {
  title: "DEPLOYMENT/DeploymentItem",
  component: DeploymentItem,
};

export const Overview: Story = () => (
  <Provider
    store={createStore({
      deployments: {
        entities: {
          [dummyDeployment.id]: dummyDeployment,
        },
        ids: [dummyDeployment.id],
      },
      applications: {
        entities: {
          [dummyApplication.id]: dummyApplication,
        },
        ids: [dummyApplication.id],
      },
    })}
  >
    <DeploymentItem id={dummyDeployment.id} />
  </Provider>
);

export const noDescription: Story = () => (
  <Provider
    store={createStore({
      deployments: {
        entities: {
          [dummyDeployment.id]: {
            ...dummyDeployment,
            summary: "",
          },
        },
        ids: [dummyDeployment.id],
      },
      applications: {
        entities: {
          [dummyApplication.id]: dummyApplication,
        },
        ids: [dummyApplication.id],
      },
    })}
  >
    <DeploymentItem id={dummyDeployment.id} />
  </Provider>
);
