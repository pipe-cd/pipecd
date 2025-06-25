import userEvent from "@testing-library/user-event";
import { LogSeverity } from "~~/model/logblock_pb";
import { dummyDeployment } from "~/__fixtures__/dummy-deployment";
import { render, screen } from "~~/test-utils";
import { LogViewer } from ".";

Element.prototype.scrollIntoView = jest.fn();

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

const activeStage = dummyDeployment.stagesList[0];

it("should not appear in the document if activeState is null", () => {
  const changeActiveStage = jest.fn();
  render(
    <LogViewer
      activeStage={null}
      stageLog={dummyLog}
      changeActiveStage={changeActiveStage}
    />
  );

  expect(screen.queryByTestId("log-viewer")).not.toBeInTheDocument();
});

it("should appear stage log in the document if activeState is exists", () => {
  render(
    <LogViewer
      activeStage={activeStage}
      stageLog={dummyLog}
      changeActiveStage={() => {}}
    />,
    {}
  );

  expect(screen.queryByTestId("log-viewer")).toBeInTheDocument();
  expect(
    screen.queryByText(dummyDeployment.stagesList[0].name)
  ).toBeInTheDocument();
  expect(screen.queryByText("hello world")).toBeInTheDocument();
});

it("should clearActiveStage action if click `close log` button", () => {
  const changeActiveStage = jest.fn();
  render(
    <LogViewer
      activeStage={activeStage}
      stageLog={dummyLog}
      changeActiveStage={changeActiveStage}
    />
  );

  userEvent.click(screen.getByRole("button", { name: "close log" }));

  expect(changeActiveStage).toHaveBeenCalled();
});
