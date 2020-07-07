import React from "react";
import { Toasts } from "./toasts";
import { createDecoratorRedux } from "../../.storybook/redux-decorator";

export default {
  title: "COMMON|Toasts",
  component: Toasts,
  decorators: [
    createDecoratorRedux({
      toasts: {
        entities: {
          "1": {
            id: "1",
            message: "toast message",
            severity: "error",
          },
        },
        ids: ["1"],
      },
    }),
  ],
};

export const overview: React.FC = () => <Toasts />;
