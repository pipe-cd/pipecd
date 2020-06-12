import React from "react";
import { Pipeline } from "./pipeline";
import {
  StageStatus,
  PipelineStage,
} from "pipe/pkg/app/web/model/deployment_pb";

export default {
  title: "Pipeline",
  component: Pipeline,
};

const makeStage = (
  id: string,
  name: string,
  requires?: string[]
): PipelineStage.AsObject => ({
  id,
  name,
  desc: "blah",
  index: 0,
  predefined: false,
  requiresList: requires || [],
  status: StageStatus.STAGE_NOT_STARTED_YET,
  metadataMap: [],
  visible: true,
  retriedCount: 1,
  completedAt: 0,
  createdAt: 0,
  updatedAt: 0,
});

const stages = [
  [makeStage("1", "K8S_PRIMARY_UPDATE")],
  [
    makeStage("2", "K8S_CANARY_ROLLOUT", ["1"]),
    makeStage("3", "K8S_CANARY_CLEAN", ["1"]),
  ],
  [makeStage("4", "K8S_CANARY_CLEAN", ["2"])],
];

export const overview: React.FC = () => <Pipeline stages={stages} />;
