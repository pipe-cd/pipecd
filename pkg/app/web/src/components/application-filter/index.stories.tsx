import * as React from "react";
import { ApplicationFilter } from "./";
import { action } from "@storybook/addon-actions";
import { createDecoratorRedux } from "../../../.storybook/redux-decorator";

export default {
  title: "APPLICATION/ApplicationFilter",
  component: ApplicationFilter,
  decorators: [
    createDecoratorRedux({
      environments: {
        entities: {
          "env-1": {
            createdAt: 0,
            desc: "env-1",
            id: "env-1",
            name: "stg",
            projectId: "1",
            updatedAt: 0,
            deletedAt: 0,
            deleted: false,
            disabled: false,
          },
        },
        ids: ["env-1"],
      },
    }),
  ],
};

export const overview: React.FC = () => (
  <ApplicationFilter
    options={{}}
    onChange={action("onChange")}
    onClear={action("onClear")}
  />
);
