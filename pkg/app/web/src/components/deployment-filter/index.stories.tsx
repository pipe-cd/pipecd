import React from "react";
import { DeploymentFilter } from "./";
import { createDecoratorRedux } from "../../../.storybook/redux-decorator";
import { dummyEnv } from "../../__fixtures__/dummy-environment";
import { dummyApplication } from "../../__fixtures__/dummy-application";
import { action } from "@storybook/addon-actions";

export default {
  title: "DEPLOYMENT/DeploymentFilter",
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
          ["test"]: { ...dummyApplication, id: "test", name: "test-app" },
        },
        ids: [dummyApplication.id, "test"],
      },
    }),
  ],
};

export const overview: React.FC = () => (
  <DeploymentFilter
    onChange={action("onChange")}
    onClear={action("onClear")}
    options={{
      applicationId: undefined,
      envId: undefined,
      kind: undefined,
      status: undefined,
    }}
  />
);
