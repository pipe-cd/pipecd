import { StageStatus } from "~~/model/deployment_pb";

export const isStageRunning = (status: StageStatus): boolean => {
  switch (status) {
    case StageStatus.STAGE_NOT_STARTED_YET:
    case StageStatus.STAGE_RUNNING:
      return true;
    case StageStatus.STAGE_SUCCESS:
    case StageStatus.STAGE_FAILURE:
    case StageStatus.STAGE_CANCELLED:
    case StageStatus.STAGE_SKIPPED:
    case StageStatus.STAGE_EXITED:
      return false;
  }
};
