import React from "react";
import { EnvironmentListItem } from "./";
import { createDecoratorRedux } from "../../../.storybook/redux-decorator";
import { dummyEnv } from "../../__fixtures__/dummy-environment";

export default {
  title: "EnvironmentListItem",
  component: EnvironmentListItem,
  decorators: [
    createDecoratorRedux({
      environments: {
        entities: { [dummyEnv.id]: dummyEnv },
        ids: [dummyEnv.id],
      },
    }),
  ],
};

export const overview: React.FC = () => (
  <ul style={{ listStyle: "none" }}>
    <EnvironmentListItem id={dummyEnv.id} />
  </ul>
);
