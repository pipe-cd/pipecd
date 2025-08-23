import { LogBlock } from "~~/model/logblock_pb";

export type StageLog = {
  deploymentId: string;
  stageId: string;
  logBlocks: LogBlock.AsObject[];
};
