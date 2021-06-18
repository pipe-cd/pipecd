import { CopyIconButton, CopyIconButtonProps } from "./";
import { Story } from "@storybook/react";
import { createDecoratorRedux } from "~~/.storybook/redux-decorator";

export default {
  title: "COMMON/CopyIconButton",
  component: CopyIconButton,
  decorators: [createDecoratorRedux({})],
};

const Template: Story<CopyIconButtonProps> = (args) => (
  <CopyIconButton {...args} />
);

export const Overview = Template.bind({});
Overview.args = {
  name: "ID",
  value: "id",
};
