import React from "react";
import { ApplicationCount, ApplicationCountProps } from "./";
import { Story } from "@storybook/react/types-6-0";
import { ApplicationKind } from "../../modules/applications";

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
  kind: ApplicationKind.KUBERNETES,
};
