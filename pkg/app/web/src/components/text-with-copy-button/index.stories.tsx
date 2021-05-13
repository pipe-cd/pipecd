import { Story } from "@storybook/react/types-6-0";
import { createDecoratorRedux } from "../../../.storybook/redux-decorator";
import { TextWithCopyButton, TextWithCopyButtonProps } from "./";

export default {
  title: "COMMON/TextWithCopyButton",
  component: TextWithCopyButton,
  decorators: [createDecoratorRedux({})],
};

const Template: Story<TextWithCopyButtonProps> = (args) => (
  <TextWithCopyButton {...args} />
);

export const Overview = Template.bind({});
Overview.args = { name: "Value", value: "value" };
