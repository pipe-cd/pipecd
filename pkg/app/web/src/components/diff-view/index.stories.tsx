import { Story } from "@storybook/react";
import { DiffView, DiffViewProps } from "./";

export default {
  title: "DiffView",
  component: DiffView,
};

const content = `
+ added line
- deleted line
normal
  indent
`;

const Template: Story<DiffViewProps> = (args) => <DiffView {...args} />;
export const Overview = Template.bind({});
Overview.args = { content };
