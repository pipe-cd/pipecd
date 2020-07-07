import React from "react";
import { createDecoratorRedux } from "../../.storybook/redux-decorator";
import { dummyDeployment } from "../__fixtures__/dummy-deployment";
import { DeploymentDetail } from "./deployment-detail";
import { dummyEnv } from "../__fixtures__/dummy-environment";

export default {
  title: "DEPLOYMENT|DeploymentDetail",
  component: DeploymentDetail,
  decorators: [
    createDecoratorRedux({
      deployments: {
        canceling: {},
        entities: {
          "deployment-1": dummyDeployment,
        },
        ids: ["deployment-1"],
        loading: {},
      },
      environments: {
        entities: {
          "env-1": dummyEnv,
        },
        ids: ["env-1"],
      },
    }),
  ],
};

export const overview: React.FC = () => (
  <DeploymentDetail deploymentId="deployment-1" />
);
