import { SplitButton, SplitButtonProps } from "./";
import CancelIcon from "@material-ui/icons/Cancel";
import { Story } from "@storybook/react";

export default {
  title: "COMMON/SplitButton",
  component: SplitButton,
  argTypes: {
    onClick: { action: "onClick" },
  },
};

const Template: Story<SplitButtonProps> = (args) => <SplitButton {...args} />;
export const Overview = Template.bind({});
Overview.args = {
  label: "select option",
  startIcon: <CancelIcon />,
  options: ["Cancel", "Cancel Without Rollback"],
  disabled: false,
  loading: false,
};

export const Loading = Template.bind({});
Loading.args = {
  label: "select option",
  startIcon: <CancelIcon />,
  options: ["Cancel", "Cancel Without Rollback"],
  disabled: true,
  loading: true,
};
