import { Story } from "@storybook/react";
import { createDecoratorRedux } from "../../../.storybook/redux-decorator";
import { StaticAdminForm } from "./";

export default {
  title: "Setting/StaticAdminForm",
  component: StaticAdminForm,
  decorators: [
    createDecoratorRedux({
      project: {
        staticAdminDisabled: false,
        username: "pipe-user",
      },
    }),
  ],
};

const Template: Story = (args) => <StaticAdminForm {...args} />;
export const Overview = Template.bind({});
Overview.args = {};
