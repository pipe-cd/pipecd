import { action } from "@storybook/addon-actions";
import React from "react";
import { createDecoratorRedux } from "../../../.storybook/redux-decorator";
import { dummyEnv } from "../../__fixtures__/dummy-environment";
import { dummyPiped } from "../../__fixtures__/dummy-piped";
import { EditPipedDrawer } from "./";

const env2 = { ...dummyEnv, id: "env-2", name: "development" };

export default {
  title: "SETTINGS/Piped/EditPipedDrawer",
  component: EditPipedDrawer,
  decorators: [
    createDecoratorRedux({
      pipeds: {
        entities: {
          [dummyPiped.id]: dummyPiped,
        },
        ids: [dummyPiped.id],
      },
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
  <EditPipedDrawer pipedId={dummyPiped.id} onClose={action("onClose")} />
);
