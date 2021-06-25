import { Story } from "@storybook/react";
import { ResourceFilterPopover, ResourceFilterPopoverProps } from ".";

export default {
  title: "ResourceFilterPopover",
  component: ResourceFilterPopover,
  argTypes: {
    onChange: {
      action: "onChange",
    },
  },
};

const Template: Story<ResourceFilterPopoverProps> = (args) => (
  <ResourceFilterPopover {...args} />
);
export const Overview = Template.bind({});
Overview.args = { enables: { Pod: true, ReplicaSet: false } };
