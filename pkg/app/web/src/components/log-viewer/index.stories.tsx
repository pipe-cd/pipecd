import * as React from "react";
import { LogViewer } from "./";
import { createDecoratorRedux } from "../../../.storybook/redux-decorator";
import { LogSeverity, createActiveStageKey } from "../../modules/stage-logs";
import { dummyDeployment } from "../../__fixtures__/dummy-deployment";
import { dummyPipelineStage } from "../../__fixtures__/dummy-pipeline";

export default {
  title: "DEPLOYMENT/LogViewer",
  component: LogViewer,
  decorators: [
    createDecoratorRedux({
      activeStage: {
        deploymentId: dummyDeployment.id,
        stageId: dummyPipelineStage.id,
        name: "active-log",
      },
      deployments: {
        entities: {
          [dummyDeployment.id]: dummyDeployment,
        },
        ids: [dummyDeployment.id],
      },
      stageLogs: {
        [createActiveStageKey({
          deploymentId: dummyDeployment.id,
          stageId: dummyPipelineStage.id,
        })]: {
          completed: true,
          deploymentId: dummyDeployment.id,
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
          stageId: dummyPipelineStage.id,
        },
      },
    }),
  ],
};

export const overview: React.FC = () => (
  <div style={{ position: "relative", height: "100vh" }}>
    <LogViewer />
  </div>
);
