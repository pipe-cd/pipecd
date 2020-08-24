import React from "react";
import { ApplicationStateView } from "./application-state-view";
import { createDecoratorRedux } from "../../.storybook/redux-decorator";
import { dummyApplicationLiveState } from "../__fixtures__/dummy-application-live-state";

export default {
  title: "APPLICATION/ApplicationStateView",
  component: ApplicationStateView,
  decorators: [
    createDecoratorRedux({
      applicationLiveState: {
        entities: {
          [dummyApplicationLiveState.applicationId]: dummyApplicationLiveState,
        },
        ids: [dummyApplicationLiveState.applicationId],
      },
    }),
  ],
};

export const overview: React.FC = () => (
  <ApplicationStateView applicationId="application-1" />
);
