import { Story } from "@storybook/react";
import { WeekPicker, WeekPickerProps } from "./";

export default {
  title: "COMMON/WeekPicker",
  component: WeekPicker,
  argTypes: {
    onChange: { action: "onChange" },
  },
};

const Template: Story<WeekPickerProps> = (args) => <WeekPicker {...args} />;
export const Overview = Template.bind({});
Overview.args = {
  value: new Date(),
  label: "Week Picker",
};
