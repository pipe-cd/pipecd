import { Story } from "@storybook/react";
import { dummyAPIKey } from "~/__fixtures__/dummy-api-key";
import { createDecoratorRedux } from "~~/.storybook/redux-decorator";
import { DisableAPIKeyConfirmDialog, DisableAPIKeyConfirmDialogProps } from ".";

export default {
  title: "Setting/APIKey/DisableAPIKeyConfirmDialog",
  component: DisableAPIKeyConfirmDialog,
  decorators: [
    createDecoratorRedux({
      apiKeys: {
        entities: { [dummyAPIKey.id]: dummyAPIKey },
        ids: [dummyAPIKey.id],
      },
    }),
  ],
  argTypes: {
    onCancel: {
      action: "onCancel",
    },
    onDisable: {
      action: "onDisable",
    },
  },
};

const Template: Story<DisableAPIKeyConfirmDialogProps> = (args) => (
  <DisableAPIKeyConfirmDialog {...args} />
);
export const Overview = Template.bind({});
Overview.args = {
  apiKeyId: dummyAPIKey.id,
};
