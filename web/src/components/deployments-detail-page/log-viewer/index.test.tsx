import userEvent from "@testing-library/user-event";
import { clearActiveStage } from "~/modules/active-stage";
import { createActiveStageKey, LogSeverity } from "~/modules/stage-logs";
import { dummyDeployment } from "~/__fixtures__/dummy-deployment";
import { createStore, render, screen } from "~~/test-utils";
import { LogViewer } from ".";

const activeStageId = createActiveStageKey({
  deploymentId: dummyDeployment.id,
  stageId: dummyDeployment.stagesList[0].id,
});

const dummyLog = {
  stageId: dummyDeployment.stagesList[0].id,
  logBlocks: [
    {
      createdAt: 0,
      index: 0,
      log: "hello world",
      severity: LogSeverity.SUCCESS,
    },
  ],
  completed: true,
  deploymentId: dummyDeployment.id,
};

it("should not appear in the document if activeState is null", () => {
  render(<LogViewer />, {
    initialState: {
      deployments: {
        ids: [dummyDeployment.id],
        entities: {
          [dummyDeployment.id]: dummyDeployment,
        },
      },
      activeStage: null,
    },
  });

  expect(screen.queryByTestId("log-viewer")).not.toBeInTheDocument();
});

it("should appear stage log in the document if activeState is exists", () => {
  render(<LogViewer />, {
    initialState: {
      deployments: {
        ids: [dummyDeployment.id],
        entities: {
          [dummyDeployment.id]: dummyDeployment,
        },
        skippable: {},
      },
      activeStage: {
        deploymentId: dummyDeployment.id,
        name: dummyDeployment.stagesList[0].name,
        stageId: dummyDeployment.stagesList[0].id,
      },
      stageLogs: {
        [activeStageId]: dummyLog,
      },
    },
  });

  expect(screen.queryByTestId("log-viewer")).toBeInTheDocument();
  expect(
    screen.queryByText(dummyDeployment.stagesList[0].name)
  ).toBeInTheDocument();
  expect(screen.queryByText("hello world")).toBeInTheDocument();
});

it("should dispatch clearActiveStage action if click `close log` button", () => {
  const store = createStore({
    deployments: {
      ids: [dummyDeployment.id],
      entities: {
        [dummyDeployment.id]: dummyDeployment,
      },
      skippable: {},
    },
    activeStage: {
      deploymentId: dummyDeployment.id,
      name: dummyDeployment.stagesList[0].name,
      stageId: dummyDeployment.stagesList[0].id,
    },
    stageLogs: {
      [activeStageId]: dummyLog,
    },
  });
  render(<LogViewer />, {
    store,
  });

  userEvent.click(screen.getByRole("button", { name: "close log" }));

  expect(store.getActions()).toMatchObject([
    {
      type: clearActiveStage.type,
    },
  ]);
});
