import React from "react";
import { ApplicationFormDrawer } from "./application-form-drawer";
import { action } from "@storybook/addon-actions";
import { createDecoratorRedux } from "../../.storybook/redux-decorator";
import { dummyEnv } from "../__fixtures__/dummy-environment";
import { dummyPiped } from "../__fixtures__/dummy-piped";

export default {
  title: "APPLICATION/ApplicationFormDrawer",
  component: ApplicationFormDrawer,
  decorators: [
    createDecoratorRedux({
      environments: {
        entities: {
          [dummyEnv.id]: dummyEnv,
        },
        ids: [dummyEnv.id],
      },
      pipeds: {
        entities: {
          [dummyPiped.id]: dummyPiped,
        },
        ids: [dummyPiped.id],
      },
    }),
  ],
};

export const overview: React.FC = () => (
  <ApplicationFormDrawer
    open
    title="add application"
    onSubmit={action("onSubmit")}
    onClose={action("onClose")}
    isProcessing={false}
  />
);

export const isProcessing: React.FC = () => (
  <ApplicationFormDrawer
    open
    title="add application"
    onSubmit={action("onSubmit")}
    onClose={action("onClose")}
    isProcessing={true}
  />
);
