import { Story } from "@storybook/react";
import { createDecoratorRedux } from "~~/.storybook/redux-decorator";
import { GithubSSOForm } from "./";

export default {
  title: "Setting/GithubSSOForm",
  component: GithubSSOForm,
  decorators: [createDecoratorRedux({})],
};

const Template: Story = (args) => <GithubSSOForm {...args} />;
export const Overview = Template.bind({});
Overview.args = {};
