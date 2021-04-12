import * as React from "react";
import { createDecoratorRedux } from "../../../.storybook/redux-decorator";
import { dummyApplication } from "../../__fixtures__/dummy-application";
import { DeleteApplicationDialog } from "./";

export default {
  title: "DeleteApplicationDialog",
  component: DeleteApplicationDialog,
  decorators: [
    createDecoratorRedux({
      applications: {
        entities: {
          [dummyApplication.id]: dummyApplication,
        },
        ids: [dummyApplication.id],
      },
      deleteApplication: {
        applicationId: dummyApplication.id,
        deleting: false,
      },
    }),
  ],
};

export const overview: React.FC = () => <DeleteApplicationDialog />;
