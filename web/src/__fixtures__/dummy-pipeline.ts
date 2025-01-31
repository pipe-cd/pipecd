import * as jspb from "google-protobuf";
import { PipelineStage, StageStatus } from "~/modules/deployments";
import { ManualOperation } from "~~/model/deployment_pb";
import { createRandTimes, randomUUID, randomWords } from "./utils";

const [createdAt, updatedAt, completedAt] = createRandTimes(3);

export const dummyPipelineStage: PipelineStage.AsObject = {
  id: randomUUID(),
  name: "K8S_CANARY_ROLLOUT",
  desc: randomWords(8),
  index: 0,
  predefined: true,
  requiresList: [],
  visible: true,
  status: StageStatus.STAGE_SUCCESS,
  statusReason: "",
  metadataMap: [],
  retriedCount: 0,
  rollback: false,
  completedAt: completedAt.unix(),
  createdAt: createdAt.unix(),
  updatedAt: updatedAt.unix(),
  availableOperation: ManualOperation.MANUAL_OPERATION_UNKNOWN,
};

export function createPipelineStage(
  pipeline: Partial<PipelineStage.AsObject>
): PipelineStage.AsObject {
  return Object.assign({}, dummyPipelineStage, pipeline, {
    id: randomUUID(),
  });
}

export const dummyPipeline: PipelineStage.AsObject[] = [dummyPipelineStage];
dummyPipeline.push(
  createPipelineStage({
    index: 1,
    name: "K8S_TRAFFIC_ROUTING",
    requiresList: [dummyPipeline[0].id],
    status: StageStatus.STAGE_RUNNING,
    metadataMap: [
      ["primary-percentage", "50"],
      ["canary-percentage", "50"],
      ["baseline-percentage", "0"],
    ],
  })
);
dummyPipeline.push(
  createPipelineStage({
    index: 2,
    name: "K8S_CANARY_CLEAN",
    requiresList: [dummyPipeline[1].id],
    status: StageStatus.STAGE_NOT_STARTED_YET,
  })
);

export function createPipelineStageFromObject(
  o: PipelineStage.AsObject
): PipelineStage {
  const stage = new PipelineStage();
  stage.setId(o.id);
  stage.setName(o.name);
  stage.setDesc(o.desc);
  stage.setIndex(o.index);
  stage.setPredefined(o.predefined);
  stage.setRequiresList(o.requiresList);
  stage.setVisible(o.visible);
  stage.setStatus(o.status);
  stage.setStatusReason(o.statusReason);
  stage.setRetriedCount(o.retriedCount);
  stage.setCompletedAt(o.completedAt);
  stage.setCreatedAt(o.createdAt);
  stage.setUpdatedAt(o.updatedAt);
  const metadataMap: jspb.Map<string, string> = stage.getMetadataMap();
  o.metadataMap.forEach((value) => {
    metadataMap.set(value[0], value[1]);
  });
  return stage;
}

export function createPipelineFromObject(
  o: PipelineStage.AsObject[]
): PipelineStage[] {
  return o.map(createPipelineStageFromObject);
}
