import { Story } from "@storybook/react";
import { ApplicationCount, ApplicationCountProps } from ".";

export default {
  title: "application/ApplicationCount",
  component: ApplicationCount,
  argTypes: {
    onClick: {
      action: "onClick",
    },
  },
};

const Template: Story<ApplicationCountProps> = (args) => (
  <ApplicationCount {...args} />
);

export const Overview = Template.bind({});
Overview.args = {
  enabledCount: 123,
  disabledCount: 12,
  kindName: "KUBERNETES",
};
