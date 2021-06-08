import { Story } from "@storybook/react";
import { createDecoratorRedux } from "../../../.storybook/redux-decorator";
import { dummyEnv } from "../../__fixtures__/dummy-environment";
import { dummyPiped } from "../../__fixtures__/dummy-piped";
import { EditPipedDrawer, EditPipedDrawerProps } from "./";

const env2 = { ...dummyEnv, id: "env-2", name: "development" };

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

const Template: Story<EditPipedDrawerProps> = (args) => (
  <EditPipedDrawer {...args} />
);
export const Overview = Template.bind({});
Overview.args = { pipedId: dummyPiped.id };
