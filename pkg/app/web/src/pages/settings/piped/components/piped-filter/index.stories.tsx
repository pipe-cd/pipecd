import { Story } from "@storybook/react";
import { PipedFilter, PipedFilterProps } from ".";

export default {
  title: "Setting/Piped/PipedFilter",
  component: PipedFilter,
  argTypes: {
    onChange: {
      action: "onChange",
    },
  },
};

const Template: Story<PipedFilterProps> = (args) => <PipedFilter {...args} />;
export const Overview = Template.bind({});
Overview.args = {
  values: {
    enabled: false,
  },
};
