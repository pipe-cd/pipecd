import { Story } from "@storybook/react";
import { AddEnvForm, AddEnvFormProps } from "./";

export default {
  title: "Setting/AddEnvForm",
  component: AddEnvForm,
  argTypes: {
    onCancel: {
      action: "onCancel",
    },
    onSubmit: {
      action: "onSubmit",
    },
  },
};

const Template: Story<AddEnvFormProps> = (args) => <AddEnvForm {...args} />;

export const Overview = Template.bind({});
Overview.args = {
  projectName: "project-name",
};
