import React from "react";
import { DeploymentItem } from "./";
import { dummyDeployment } from "../../__fixtures__/dummy-deployment";
import { dummyApplication } from "../../__fixtures__/dummy-application";
import { dummyEnv } from "../../__fixtures__/dummy-environment";
import { Provider } from "react-redux";
import { createStore } from "../../../test-utils";

export default {
  title: "DEPLOYMENT/DeploymentItem",
  component: DeploymentItem,
};

export const overview: React.FC = () => (
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

export const noDescription: React.FC = () => (
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
