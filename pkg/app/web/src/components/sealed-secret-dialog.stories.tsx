import { action } from "@storybook/addon-actions";
import React from "react";
import { createDecoratorRedux } from "../../.storybook/redux-decorator";
import { dummyApplication } from "../__fixtures__/dummy-application";
import { SealedSecretDialog } from "./sealed-secret-dialog";

export default {
  title: "SealedSecretDialog",
  component: SealedSecretDialog,
  decorators: [
    createDecoratorRedux({
      applications: {
        entities: { [dummyApplication.id]: dummyApplication },
        ids: [dummyApplication.id],
      },
    }),
  ],
};

export const overview: React.FC = () => (
  <SealedSecretDialog
    open
    applicationId={dummyApplication.id}
    onClose={action("onClose")}
  />
);
