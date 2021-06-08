import { Story } from "@storybook/react";
import { createDecoratorRedux } from "../../../.storybook/redux-decorator";
import { GeneratedAPIKeyDialog } from "./";

export default {
  title: "Setting/APIKey/GeneratedAPIKeyDialog",
  component: GeneratedAPIKeyDialog,
  decorators: [
    createDecoratorRedux({
      apiKeys: {
        disabling: false,
        error: null,
        generatedKey:
          "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.bspmf2xvt74area19iaxl0yh33jzwelq493vzil0orgzylrdb1",
        generating: false,
        loading: false,
        entities: {},
        ids: [],
      },
    }),
  ],
};

const Template: Story = (args) => <GeneratedAPIKeyDialog {...args} />;
export const Overview = Template.bind({});
Overview.args = {};
