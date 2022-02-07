import { Story } from "@storybook/react";
import { DeploymentStatus } from "~/modules/deployments";
import { dummyDeployment } from "~/__fixtures__/dummy-deployment";
import { dummyPiped } from "~/__fixtures__/dummy-piped";
import { createDecoratorRedux } from "~~/.storybook/redux-decorator";
import { DeploymentDetail, DeploymentDetailProps } from ".";

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
