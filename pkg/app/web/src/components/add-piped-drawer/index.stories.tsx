import { Story } from "@storybook/react";
import { createDecoratorRedux } from "../../../.storybook/redux-decorator";
import { dummyEnv } from "../../__fixtures__/dummy-environment";
import { AddPipedDrawer, AddPipedDrawerProps } from "./";

const env2 = { ...dummyEnv, id: "env-2", name: "development" };

export default {
  title: "Setting/Piped/AddPipedDrawer",
  component: AddPipedDrawer,
  decorators: [
    createDecoratorRedux({
      environments: {
        entities: {
          [dummyEnv.id]: dummyEnv,
          [env2.id]: env2,
        },
        ids: [dummyEnv.id, env2.id],
      },
    }),
  ],
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
