import { Story } from "@storybook/react";
import { GenerateAPIKeyDialog, GenerateAPIKeyDialogProps } from "./";

export default {
  title: "Setting/APIKey/GenerateAPIKeyDialog",
  component: GenerateAPIKeyDialog,
  argTypes: {
    onClose: {
      action: "onClose",
    },
    onSubmit: {
      action: "onSubmit",
    },
  },
};

const Template: Story<GenerateAPIKeyDialogProps> = (args) => (
  <GenerateAPIKeyDialog {...args} />
);
export const Overview = Template.bind({});
Overview.args = { open: true };
