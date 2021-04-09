import * as React from "react";
import { createDecoratorRedux } from "../../../.storybook/redux-decorator";
import { StaticAdminForm } from "./";

export default {
  title: "SETTINGS/StaticAdminForm",
  component: StaticAdminForm,
  decorators: [
    createDecoratorRedux({
      project: {
        staticAdminDisabled: false,
        username: "pipe-user",
      },
    }),
  ],
};

export const overview: React.FC = () => <StaticAdminForm />;
