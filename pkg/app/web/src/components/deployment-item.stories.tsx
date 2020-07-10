import React from "react";
import { DeploymentItem } from "./deployment-item";
import { createDecoratorRedux } from "../../.storybook/redux-decorator";
import { dummyDeployment } from "../__fixtures__/dummy-deployment";
import { dummyApplication } from "../__fixtures__/dummy-application";
import { dummyEnv } from "../__fixtures__/dummy-environment";

export default {
  title: "DEPLOYMENT|DeploymentItem",
  component: DeploymentItem,
  decorators: [
    createDecoratorRedux({
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
    }),
  ],
};

export const overview: React.FC = () => (
  <DeploymentItem id={dummyDeployment.id} />
);
