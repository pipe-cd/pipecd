import { dummyDeployment } from "~/__fixtures__/dummy-deployment";
import { createPipelineStage } from "~/__fixtures__/dummy-pipeline";
import { Deployment, StageStatus } from "~~/model/deployment_pb";
import { findDefaultActiveStageInDeployment } from "./find-default-active-stage-in-deployment";

describe("findDefaultActiveStageInDeployment", () => {
  it("returns null when deployment is undefined", () => {
    expect(findDefaultActiveStageInDeployment(undefined)).toBeNull();
  });

  describe("Piped v0", () => {
    it("returns the running stage if one exists", () => {
      const stageSuccess = createPipelineStage({
        index: 0,
        visible: true,
        status: StageStatus.STAGE_SUCCESS,
      });
      const stageRunning = createPipelineStage({
        index: 1,
        visible: true,
        status: StageStatus.STAGE_RUNNING,
      });
      const stageNotStarted = createPipelineStage({
        index: 2,
        visible: true,
        status: StageStatus.STAGE_NOT_STARTED_YET,
      });

      const deployment = {
        ...dummyDeployment,
        deployTargetsByPluginMap: [],
        stagesList: [stageSuccess, stageRunning, stageNotStarted],
      };

      expect(findDefaultActiveStageInDeployment(deployment)).toEqual(
        stageRunning
      );
    });

    it("returns the latest visible active stage when none are running", () => {
      const stage1 = createPipelineStage({
        index: 0,
        visible: true,
        status: StageStatus.STAGE_SUCCESS,
      });
      const stage2 = createPipelineStage({
        index: 1,
        visible: true,
        status: StageStatus.STAGE_SUCCESS,
      });
      const stage3 = createPipelineStage({
        index: 2,
        visible: false,
        status: StageStatus.STAGE_SUCCESS,
      });

      const deployment = {
        ...dummyDeployment,
        deployTargetsByPluginMap: [],
        stagesList: [stage1, stage2, stage3],
      };

      expect(findDefaultActiveStageInDeployment(deployment)).toEqual(stage2);
    });

    it("filters out not-started stages", () => {
      const stage1 = createPipelineStage({
        index: 0,
        visible: true,
        status: StageStatus.STAGE_SUCCESS,
      });
      const stage2 = createPipelineStage({
        index: 1,
        visible: true,
        status: StageStatus.STAGE_NOT_STARTED_YET,
      });

      const deployment = {
        ...dummyDeployment,
        deployTargetsByPluginMap: [],
        stagesList: [stage1, stage2],
      };

      expect(findDefaultActiveStageInDeployment(deployment)).toEqual(stage1);
    });
    it("returns null when all stages are filtered out (e.g. not started)", () => {
      const stage1 = createPipelineStage({
        index: 0,
        visible: true,
        status: StageStatus.STAGE_NOT_STARTED_YET,
      });

      const deployment = {
        ...dummyDeployment,
        deployTargetsByPluginMap: [],
        stagesList: [stage1],
      };

      expect(findDefaultActiveStageInDeployment(deployment)).toBeNull();
    });

    it("returns null when all stages are invisible", () => {
      const stage1 = createPipelineStage({
        index: 0,
        visible: false,
        status: StageStatus.STAGE_SUCCESS,
      });

      const deployment = {
        ...dummyDeployment,
        deployTargetsByPluginMap: [],
        stagesList: [stage1],
      };

      expect(findDefaultActiveStageInDeployment(deployment)).toBeNull();
    });
  });

  describe("Piped v1", () => {
    it("ignores stage.visible and returns running stage", () => {
      const stageSuccess = createPipelineStage({
        index: 0,
        visible: false,
        status: StageStatus.STAGE_SUCCESS,
      });
      const stageRunning = createPipelineStage({
        index: 1,
        visible: false,
        status: StageStatus.STAGE_RUNNING,
      });

      const deployment = {
        ...dummyDeployment,
        deployTargetsByPluginMap: [
          ["plugin-1", { deployTargetsList: [] }],
        ] as Deployment.AsObject["deployTargetsByPluginMap"],
        stagesList: [stageSuccess, stageRunning],
      };

      expect(findDefaultActiveStageInDeployment(deployment)).toEqual(
        stageRunning
      );
    });

    it("returns latest non-not-started stage when none are running", () => {
      const stage1 = createPipelineStage({
        index: 0,
        visible: false,
        status: StageStatus.STAGE_SUCCESS,
      });
      const stage2 = createPipelineStage({
        index: 1,
        visible: false,
        status: StageStatus.STAGE_SUCCESS,
      });
      const stage3 = createPipelineStage({
        index: 2,
        visible: false,
        status: StageStatus.STAGE_NOT_STARTED_YET,
      });

      const deployment = {
        ...dummyDeployment,
        deployTargetsByPluginMap: [
          ["plugin-1", { deployTargetsList: [] }],
        ] as Deployment.AsObject["deployTargetsByPluginMap"],
        stagesList: [stage1, stage2, stage3],
      };

      expect(findDefaultActiveStageInDeployment(deployment)).toEqual(stage2);
    });

    it("returns null when all stages are not started", () => {
      const stage1 = createPipelineStage({
        index: 0,
        visible: false,
        status: StageStatus.STAGE_NOT_STARTED_YET,
      });

      const deployment = {
        ...dummyDeployment,
        deployTargetsByPluginMap: [
          ["plugin-1", { deployTargetsList: [] }],
        ] as Deployment.AsObject["deployTargetsByPluginMap"],
        stagesList: [stage1],
      };

      expect(findDefaultActiveStageInDeployment(deployment)).toBeNull();
    });

    it("returns null when stages list is empty", () => {
      const deployment = {
        ...dummyDeployment,
        deployTargetsByPluginMap: [
          ["plugin-1", { deployTargetsList: [] }],
        ] as Deployment.AsObject["deployTargetsByPluginMap"],
        stagesList: [],
      };

      expect(findDefaultActiveStageInDeployment(deployment)).toBeNull();
    });
  });
});
