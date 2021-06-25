import { DeploymentFilter, DeploymentFilterProps } from ".";
import { createDecoratorRedux } from "~~/.storybook/redux-decorator";
import { dummyEnv } from "~/__fixtures__/dummy-environment";
import { dummyApplication } from "~/__fixtures__/dummy-application";
import { Story } from "@storybook/react";

export default {
  title: "DEPLOYMENT/DeploymentFilter",
  component: DeploymentFilter,
  decorators: [
    createDecoratorRedux({
      environments: {
        entities: {
          [dummyEnv.id]: dummyEnv,
        },
        ids: [dummyEnv.id],
      },
      applications: {
        entities: {
          [dummyApplication.id]: dummyApplication,
          ["test"]: { ...dummyApplication, id: "test", name: "test-app" },
        },
        ids: [dummyApplication.id, "test"],
      },
    }),
  ],
  argTypes: {
    onChange: {
      action: "onChange",
    },
    onClear: {
      action: "onClear",
    },
  },
};

const Template: Story<DeploymentFilterProps> = (args) => (
  <DeploymentFilter {...args} />
);
export const Overview = Template.bind({});
Overview.args = {
  options: {
    applicationId: undefined,
    envId: undefined,
    kind: undefined,
    status: undefined,
  },
};
