import { DeploymentFilter, DeploymentFilterProps } from ".";
import { createDecoratorRedux } from "~~/.storybook/redux-decorator";
import { dummyApplication } from "~/__fixtures__/dummy-application";
import { Story } from "@storybook/react";

export default {
  title: "DEPLOYMENT/DeploymentFilter",
  component: DeploymentFilter,
  decorators: [
    createDecoratorRedux({
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
