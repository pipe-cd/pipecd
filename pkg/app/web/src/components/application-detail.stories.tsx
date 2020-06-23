import React from "react";
import { ApplicationDetail } from "./application-detail";
import { createDecoratorRedux } from "../../.storybook/redux-decorator";

export default {
  title: "ApplicationDetail",
  component: ApplicationDetail,
  decorators: [createDecoratorRedux({})],
};

export const overview: React.FC = () => (
  <ApplicationDetail applicationId="application-1" />
);
