import { Story } from "@storybook/react";
import { createDecoratorRedux } from "~~/.storybook/redux-decorator";
import { LoginForm, LoginFormProps } from ".";

export default {
  title: "LOGIN/LoginForm",
  component: LoginForm,
  decorators: [createDecoratorRedux({})],
};

const Template: Story<LoginFormProps> = (args) => <LoginForm {...args} />;
export const Overview = Template.bind({});
Overview.args = {
  projectName: "PipeCD",
};
