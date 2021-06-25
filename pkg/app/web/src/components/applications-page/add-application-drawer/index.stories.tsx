import { Story } from "@storybook/react";
import { dummyEnv } from "~/__fixtures__/dummy-environment";
import { dummyPiped } from "~/__fixtures__/dummy-piped";
import { createDecoratorRedux } from "~~/.storybook/redux-decorator";
import { AddApplicationDrawer, AddApplicationDrawerProps } from ".";

export default {
  title: "APPLICATION/AddApplicationDrawer",
  component: AddApplicationDrawer,
  decorators: [
    createDecoratorRedux({
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
  argTypes: {
    onClose: {
      action: "onClose",
    },
  },
};

const Template: Story<AddApplicationDrawerProps> = (args) => (
  <AddApplicationDrawer {...args} />
);

export const Overview = Template.bind({});
Overview.args = {
  open: true,
};
