import { Story } from "@storybook/react";
import { FilterView, FilterViewProps } from "./";

export default {
  title: "FilterView",
  component: FilterView,
  argTypes: {
    onClear: {
      action: "onClear",
    },
  },
};

const Template: Story<FilterViewProps> = (args) => <FilterView {...args} />;
export const Overview = Template.bind({});
Overview.args = {};
