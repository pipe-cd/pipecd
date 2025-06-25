import { dummyDeployment } from "~/__fixtures__/dummy-deployment";
import { Pipeline } from ".";
import { render, screen } from "~~/test-utils";

it("should render correct stage list", () => {
  render(
    <Pipeline
      deployment={dummyDeployment}
      activeStageInfo={null}
      changeActiveStage={() => {}}
    />
  );

  expect(screen.getByText("K8S_CANARY_ROLLOUT")).toBeInTheDocument();
  expect(screen.getByText("K8S_TRAFFIC_ROUTING")).toBeInTheDocument();
  expect(screen.getByText("K8S_CANARY_CLEAN")).toBeInTheDocument();
});

it("should call setActiveStage when click stage", () => {
  const setActiveStage = jest.fn();

  render(
    <Pipeline
      deployment={dummyDeployment}
      activeStageInfo={null}
      changeActiveStage={setActiveStage}
    />
  );

  const stage = screen.getByText("K8S_CANARY_ROLLOUT");
  stage.click();

  const expectedStage = dummyDeployment.stagesList.find(
    (item) => item.name === "K8S_CANARY_ROLLOUT"
  );

  expect(setActiveStage).toHaveBeenCalledWith({
    deploymentId: dummyDeployment.id,
    stageId: expectedStage?.id,
    name: expectedStage?.name,
  });
});
