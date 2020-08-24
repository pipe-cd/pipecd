import React from "react";
import { ApplicationList } from "./application-list";
import { createDecoratorRedux } from "../../.storybook/redux-decorator";
import { dummyApplication } from "../__fixtures__/dummy-application";

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
    }),
  ],
};

export const overview: React.FC = () => <ApplicationList />;
