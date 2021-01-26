import React from "react";
import { createDecoratorRedux } from "../../../.storybook/redux-decorator";
import { dummyDeployment } from "../../__fixtures__/dummy-deployment";
import { DeploymentDetail } from "./deployment-detail";
import { dummyEnv } from "../../__fixtures__/dummy-environment";
import { dummyPiped } from "../../__fixtures__/dummy-piped";
import { DeploymentStatus } from "../../modules/deployments";

export default {
  title: "DEPLOYMENT/DeploymentDetail",
  component: DeploymentDetail,
  decorators: [
    createDecoratorRedux({
      deployments: {
        canceling: {},
        entities: {
          [dummyDeployment.id]: dummyDeployment,
          [dummyDeployment.id + 1]: {
            ...dummyDeployment,
            id: dummyDeployment.id + 1,
            status: DeploymentStatus.DEPLOYMENT_RUNNING,
          },
        },
        ids: [dummyDeployment.id],
        loading: {},
      },
      environments: {
        entities: {
          [dummyEnv.id]: dummyEnv,
        },
        ids: [dummyEnv.id],
      },
      pipeds: {
        entities: {
          [dummyPiped.id]: dummyPiped,
        },
        ids: [dummyPiped.id],
      },
    }),
  ],
};

export const overview: React.FC = () => (
  <DeploymentDetail deploymentId={dummyDeployment.id} />
);

export const Running: React.FC = () => (
  <DeploymentDetail deploymentId={dummyDeployment.id + 1} />
);
