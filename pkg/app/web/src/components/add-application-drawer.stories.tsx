import React from "react";
import { AddApplicationDrawer } from "./add-application-drawer";
import { action } from "@storybook/addon-actions";
import { createDecoratorRedux } from "../../.storybook/redux-decorator";

export default {
  title: "APPLICATION/AddApplicationDrawer",
  component: AddApplicationDrawer,
  decorators: [createDecoratorRedux({})],
};

export const overview: React.FC = () => (
  <AddApplicationDrawer
    open
    projectName="pipe-cd"
    onSubmit={action("onSubmit")}
    onClose={action("onClose")}
    isAdding={false}
  />
);
