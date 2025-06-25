import { Stage } from "~/types/deployment";
import { Deployment, StageStatus } from "~~/model/deployment_pb";

// Find stage that is running or latest
export const findDefaultActiveStageInDeployment = (
  deployment: Deployment.AsObject | undefined
): Stage | null => {
  if (!deployment) {
    return null;
  }

  const stages = deployment.stagesList.filter(
    (stage) =>
      stage.visible && stage.status !== StageStatus.STAGE_NOT_STARTED_YET
  );

  const runningStage = stages.find(
    (stage) => stage.status === StageStatus.STAGE_RUNNING
  );

  if (runningStage) {
    return runningStage;
  }

  return stages[stages.length - 1];
};
