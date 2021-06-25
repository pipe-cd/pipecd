import { Story } from "@storybook/react";
import { createDecoratorRedux } from "~~/.storybook/redux-decorator";
import { InsightHeader } from ".";

export default {
  title: "INSIGHTS/InsightHeader",
  component: InsightHeader,
  decorators: [createDecoratorRedux({})],
};

const Template: Story = (args) => <InsightHeader {...args} />;
export const Overview = Template.bind({});
Overview.args = {};
