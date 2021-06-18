import { createDecoratorRedux } from "~~/.storybook/redux-decorator";
import { dummyDeployment } from "~/__fixtures__/dummy-deployment";
import { DeploymentDetail, DeploymentDetailProps } from "./";
import { dummyEnv } from "~/__fixtures__/dummy-environment";
import { dummyPiped } from "~/__fixtures__/dummy-piped";
import { DeploymentStatus } from "~/modules/deployments";
import { Story } from "@storybook/react";

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

const Template: Story<DeploymentDetailProps> = (args) => (
  <DeploymentDetail {...args} />
);
export const Overview = Template.bind({});
Overview.args = { deploymentId: dummyDeployment.id };

export const Running = Template.bind({});
Running.args = { deploymentId: dummyDeployment.id + 1 };
