import { Story } from "@storybook/react/types-6-0";
import React from "react";
import { ApplicationCount, ApplicationCountProps } from "./";

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
  totalCount: 123,
  disabledCount: 12,
  kindName: "KUBERNETES",
};
