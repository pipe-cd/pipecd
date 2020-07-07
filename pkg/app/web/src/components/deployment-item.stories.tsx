import React from "react";
import { DeploymentItem } from "./deployment-item";
import { createDecoratorRedux } from "../../.storybook/redux-decorator";

export default {
  title: "DEPLOYMENT|DeploymentItem",
  component: DeploymentItem,
  decorators: [createDecoratorRedux({})],
};

export const overview: React.FC = () => <DeploymentItem id="app-1" />;
