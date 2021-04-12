import { action } from "@storybook/addon-actions";
import * as React from "react";
import { createDecoratorRedux } from "../../../.storybook/redux-decorator";
import { dummyAPIKey } from "../../__fixtures__/dummy-api-key";
import { DisableAPIKeyConfirmDialog } from "./";

export default {
  title: "SETTINGS/APIKey/DisableAPIKeyConfirmDialog",
  component: DisableAPIKeyConfirmDialog,
  decorators: [
    createDecoratorRedux({
      apiKeys: {
        entities: { [dummyAPIKey.id]: dummyAPIKey },
        ids: [dummyAPIKey.id],
      },
    }),
  ],
};

export const overview: React.FC = () => (
  <DisableAPIKeyConfirmDialog
    apiKeyId={dummyAPIKey.id}
    onCancel={action("onCancel")}
    onDisable={action("onDisable")}
  />
);
