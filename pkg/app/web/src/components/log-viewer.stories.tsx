import React from "react";
import { LogViewer } from "./log-viewer";
import { createDecoratorRedux } from "../../.storybook/redux-decorator";
import { LogSeverity } from "../modules/stage-logs";

export default {
  title: "DEPLOYMENT|LogViewer",
  component: LogViewer,
  decorators: [
    createDecoratorRedux({
      activeStage: {
        id: "active-log-1",
        name: "active-log",
      },
      stageLogs: {
        "active-log-1": {
          completed: true,
          deploymentId: "1",
          logBlocks: [
            {
              createdAt: 0,
              index: 0,
              log: "HELLO",
              severity: LogSeverity.SUCCESS,
            },
            {
              createdAt: 0,
              index: 1,
              log: "ERROR",
              severity: LogSeverity.ERROR,
            },
            {
              createdAt: 0,
              index: 2,
              log: "INFO",
              severity: LogSeverity.INFO,
            },
          ],
          stageId: "1",
        },
      },
    }),
  ],
};

export const overview: React.FC = () => <LogViewer />;
