import React from "react";
import { AddPipedForm } from "./add-piped-form";
import { action } from "@storybook/addon-actions";
import { createDecoratorRedux } from "../../.storybook/redux-decorator";
import { dummyEnv } from "../__fixtures__/dummy-environment";

const env2 = { ...dummyEnv, id: "env-2", name: "development" };

export default {
  title: "SETTINGS/AddPipedForm",
  component: AddPipedForm,
  decorators: [
    createDecoratorRedux({
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
  <AddPipedForm
    onClose={action("onClose")}
    onSubmit={action("onSubmit")}
    projectName="project-name"
  />
);
