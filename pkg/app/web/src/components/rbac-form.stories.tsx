import React from "react";
import { RBACForm } from "./rbac-form";
import { createDecoratorRedux } from "../../.storybook/redux-decorator";

export default {
  title: "RBACForm",
  component: RBACForm,
  decorators: [
    createDecoratorRedux({
      project: {
        teams: { admin: "admin", editor: "editor", viewer: "viewer" },
      },
    }),
  ],
};

export const overview: React.FC = () => <RBACForm />;
