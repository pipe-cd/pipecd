import { StageStatus } from "~~/model/deployment_pb";
import { isStageRunning } from "./is-stage-running";

describe("isStageRunning", () => {
  it.each([
    [StageStatus.STAGE_NOT_STARTED_YET, true],
    [StageStatus.STAGE_RUNNING, true],
    [StageStatus.STAGE_SUCCESS, false],
    [StageStatus.STAGE_FAILURE, false],
    [StageStatus.STAGE_CANCELLED, false],
    [StageStatus.STAGE_SKIPPED, false],
    [StageStatus.STAGE_EXITED, false],
  ])("given status %p, returns %p", (status, expected) => {
    expect(isStageRunning(status)).toBe(expected);
  });
});
