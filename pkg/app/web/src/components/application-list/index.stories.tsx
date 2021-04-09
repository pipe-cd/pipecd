import * as React from "react";
import { ApplicationList } from "./";
import { createDecoratorRedux } from "../../../.storybook/redux-decorator";
import { dummyApplication } from "../../__fixtures__/dummy-application";
import { dummyEnv } from "../../__fixtures__/dummy-environment";

export default {
  title: "APPLICATION/ApplicationList",
  component: ApplicationList,
  decorators: [
    createDecoratorRedux({
      applications: {
        adding: false,
        disabling: {},
        syncing: {},
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
    }),
  ],
};

export const overview: React.FC = () => <ApplicationList />;
