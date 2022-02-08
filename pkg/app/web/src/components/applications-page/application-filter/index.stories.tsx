import { Story } from "@storybook/react";
import { createDecoratorRedux } from "~~/.storybook/redux-decorator";
import { ApplicationFilter, ApplicationFilterProps } from ".";

export default {
  title: "APPLICATION/ApplicationFilter",
  component: ApplicationFilter,
  decorators: [createDecoratorRedux({})],
  argTypes: {
    onChange: { action: "onChange" },
    onClear: { action: "onClear" },
  },
};

const Template: Story<ApplicationFilterProps> = (args) => (
  <ApplicationFilter {...args} />
);
export const Overview = Template.bind({});
Overview.args = { options: {} };
