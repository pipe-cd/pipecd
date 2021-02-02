import React from "react";
import { AddApplicationDrawer } from "./";
import { action } from "@storybook/addon-actions";
import { createDecoratorRedux } from "../../../.storybook/redux-decorator";
import { dummyEnv } from "../../__fixtures__/dummy-environment";
import { dummyPiped } from "../../__fixtures__/dummy-piped";

export default {
  title: "APPLICATION/AddApplicationDrawer",
  component: AddApplicationDrawer,
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
  <AddApplicationDrawer open onClose={action("onClose")} />
);
