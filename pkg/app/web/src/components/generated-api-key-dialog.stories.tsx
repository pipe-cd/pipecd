import { action } from "@storybook/addon-actions";
import React from "react";
import { createDecoratorRedux } from "../../.storybook/redux-decorator";
import { GeneratedApiKeyDialog } from "./generated-api-key-dialog";

export default {
  title: "GeneratedApiKeyDialog",
  component: GeneratedApiKeyDialog,
  decorators: [createDecoratorRedux({})],
};

export const overview: React.FC = () => (
  <GeneratedApiKeyDialog
    open
    generatedKey="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.bspmf2xvt74area19iaxl0yh33jzwelq493vzil0orgzylrdb1"
    onClose={action("onClose")}
  />
);
