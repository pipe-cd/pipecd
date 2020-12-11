import React from "react";
import { createDecoratorRedux } from "../../.storybook/redux-decorator";
import { dummyApplication } from "../__fixtures__/dummy-application";
import { dummyEnv } from "../__fixtures__/dummy-environment";
import { dummyPiped } from "../__fixtures__/dummy-piped";
import { EditApplicationDrawer } from "./edit-application-drawer";

export default {
  title: "APPLICATION/EditApplicationDrawer",
  component: EditApplicationDrawer,
  decorators: [
    createDecoratorRedux({
      updateApplication: {
        targetId: dummyApplication.id,
      },
      applications: {
        entities: {
          [dummyApplication.id]: dummyApplication,
        },
        ids: [dummyApplication.id],
      },
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

export const overview: React.FC = () => <EditApplicationDrawer />;
