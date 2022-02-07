import { Story } from "@storybook/react";
import { createDecoratorRedux } from "~~/.storybook/redux-decorator";
import { AddPipedDrawer, AddPipedDrawerProps } from ".";

export default {
  title: "Setting/Piped/AddPipedDrawer",
  component: AddPipedDrawer,
  decorators: [createDecoratorRedux({})],
  argTypes: {
    onClose: {
      action: "onClose",
    },
  },
};

const Template: Story<AddPipedDrawerProps> = (args) => (
  <AddPipedDrawer {...args} />
);
export const Overview = Template.bind({});
Overview.args = { open: true };
