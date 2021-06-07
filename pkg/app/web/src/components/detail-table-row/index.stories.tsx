import { Story } from "@storybook/react";
import { DetailTableRow, DetailTableRowProps } from "./";

export default {
  title: "COMMON/DetailTableRow",
  component: DetailTableRow,
};

const Template: Story<DetailTableRowProps> = (args) => (
  <DetailTableRow {...args} />
);
export const Overview = Template.bind({});
Overview.args = { label: "piped", value: "hello-world" };
