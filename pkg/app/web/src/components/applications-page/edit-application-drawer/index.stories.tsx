import { Story } from "@storybook/react";
import { dummyApplication } from "~/__fixtures__/dummy-application";
import { dummyPiped } from "~/__fixtures__/dummy-piped";
import { createDecoratorRedux } from "~~/.storybook/redux-decorator";
import { EditApplicationDrawer, EditApplicationDrawerProps } from ".";

export default {
  title: "APPLICATION/EditApplicationDrawer",
  component: EditApplicationDrawer,
  decorators: [
    createDecoratorRedux({
      updateApplication: {
        targetId: dummyApplication.id,
      },
      applications: {
        entities: {
          [dummyApplication.id]: dummyApplication,
        },
        ids: [dummyApplication.id],
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
    onUpdated: {
      action: "onUpdated",
    },
  },
};

const Template: Story<EditApplicationDrawerProps> = (args) => (
  <EditApplicationDrawer {...args} />
);

export const Overview = Template.bind({});
