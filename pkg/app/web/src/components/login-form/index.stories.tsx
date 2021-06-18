import { LoginForm, LoginFormProps } from "./";
import { createDecoratorRedux } from "~~/.storybook/redux-decorator";
import { Story } from "@storybook/react";

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
