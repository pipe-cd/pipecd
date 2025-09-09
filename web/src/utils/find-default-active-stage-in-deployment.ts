import { Stage } from "~/types/deployment";
import { Deployment, StageStatus } from "~~/model/deployment_pb";

enum PIPED_VERSION {
  V0 = "v0",
  V1 = "v1",
}

// Find stage that is running or latest
export const findDefaultActiveStageInDeployment = (
  deployment: Deployment.AsObject | undefined
): Stage | null => {
  if (!deployment) {
    return null;
  }

  const pipedVersion = deployment.deployTargetsByPluginMap?.length
    ? PIPED_VERSION.V1
    : PIPED_VERSION.V0;

  const stages = deployment.stagesList.filter((stage) => {
    if (pipedVersion === PIPED_VERSION.V0) {
      return (
        stage.visible && stage.status !== StageStatus.STAGE_NOT_STARTED_YET
      );
    }

    // For piped v1, field visible is not used.
    if (pipedVersion === PIPED_VERSION.V1) {
      return stage.status !== StageStatus.STAGE_NOT_STARTED_YET;
    }
    return false;
  });

  const runningStage = stages.find(
    (stage) => stage.status === StageStatus.STAGE_RUNNING
  );

  if (runningStage) {
    return runningStage;
  }

  return stages[stages.length - 1];
};
