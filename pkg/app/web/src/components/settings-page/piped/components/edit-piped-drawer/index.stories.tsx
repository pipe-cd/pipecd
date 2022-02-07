import { Story } from "@storybook/react";
import { dummyPiped } from "~/__fixtures__/dummy-piped";
import { createDecoratorRedux } from "~~/.storybook/redux-decorator";
import { EditPipedDrawer, EditPipedDrawerProps } from ".";

export default {
  title: "Setting/Piped/EditPipedDrawer",
  component: EditPipedDrawer,
  decorators: [
    createDecoratorRedux({
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

const Template: Story<EditPipedDrawerProps> = (args) => (
  <EditPipedDrawer {...args} />
);
export const Overview = Template.bind({});
Overview.args = { pipedId: dummyPiped.id };
