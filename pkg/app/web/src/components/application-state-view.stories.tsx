import React from "react";
import { ApplicationStateView } from "./application-state-view";
import { createDecoratorRedux } from "../../.storybook/redux-decorator";

export default {
  title: "ApplicationStateView",
  component: ApplicationStateView,
  decorators: [createDecoratorRedux({})],
};

export const overview: React.FC = () => (
  <ApplicationStateView applicationId="application-1" />
);
