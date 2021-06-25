import { Story } from "@storybook/react";
import { dummyEnv } from "~/__fixtures__/dummy-environment";
import { createDecoratorRedux } from "~~/.storybook/redux-decorator";
import { EnvironmentListItem, EnvironmentListItemProps } from ".";

export default {
  title: "EnvironmentListItem",
  component: EnvironmentListItem,
  decorators: [
    createDecoratorRedux({
      environments: {
        entities: { [dummyEnv.id]: dummyEnv },
        ids: [dummyEnv.id],
      },
    }),
  ],
};

const Template: Story<EnvironmentListItemProps> = (args) => (
  <ul style={{ listStyle: "none" }}>
    <EnvironmentListItem {...args} />
  </ul>
);
export const Overview = Template.bind({});
Overview.args = { id: dummyEnv.id };
