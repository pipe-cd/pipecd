import React from "react";
import { LogViewer } from "./log-viewer";
import { createDecoratorRedux } from "../../.storybook/redux-decorator";

export default {
  title: "LogViewer",
  component: LogViewer,
  decorators: [createDecoratorRedux({})],
};

export const overview: React.FC = () => <LogViewer />;
