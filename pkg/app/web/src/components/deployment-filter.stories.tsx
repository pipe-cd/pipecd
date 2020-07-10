import React from "react";
import { DeploymentFilter } from "./deployment-filter";
import { createDecoratorRedux } from "../../.storybook/redux-decorator";
import { dummyEnv } from "../__fixtures__/dummy-environment";
import { dummyApplication } from "../__fixtures__/dummy-application";
import { action } from "@storybook/addon-actions";

export default {
  title: "DeploymentFilter",
  component: DeploymentFilter,
  decorators: [
    createDecoratorRedux({
      environments: {
        entities: {
          [dummyEnv.id]: dummyEnv,
        },
        ids: [dummyEnv.id],
      },
      applications: {
        entities: {
          [dummyApplication.id]: dummyApplication,
        },
        ids: [dummyApplication.id],
      },
    }),
  ],
};

export const overview: React.FC = () => (
  <DeploymentFilter open onChange={action("onChange")} />
);
