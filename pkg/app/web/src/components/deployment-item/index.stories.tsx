import { DeploymentItem } from "./";
import { dummyDeployment } from "../../__fixtures__/dummy-deployment";
import { dummyApplication } from "../../__fixtures__/dummy-application";
import { dummyEnv } from "../../__fixtures__/dummy-environment";
import { Provider } from "react-redux";
import { createStore } from "test-utils";
import { Story } from "@storybook/react";

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
      environments: {
        entities: {
          [dummyEnv.id]: dummyEnv,
        },
        ids: [dummyEnv.id],
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
      environments: {
        entities: {
          [dummyEnv.id]: dummyEnv,
        },
        ids: [dummyEnv.id],
      },
    })}
  >
    <DeploymentItem id={dummyDeployment.id} />
  </Provider>
);
