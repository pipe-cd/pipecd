import { action } from "@storybook/addon-actions";
import * as React from "react";
import { createDecoratorRedux } from "../../../.storybook/redux-decorator";
import { dummyEnv } from "../../__fixtures__/dummy-environment";
import { AddPipedDrawer } from "./";

const env2 = { ...dummyEnv, id: "env-2", name: "development" };

export default {
  title: "SETTINGS/Piped/AddPipedDrawer",
  component: AddPipedDrawer,
  decorators: [
    createDecoratorRedux({
      environments: {
        entities: {
          [dummyEnv.id]: dummyEnv,
          [env2.id]: env2,
        },
        ids: [dummyEnv.id, env2.id],
      },
    }),
  ],
};

export const overview: React.FC = () => (
  <AddPipedDrawer open onClose={action("onClose")} />
);
