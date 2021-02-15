import { LogBlock, LogSeverity } from "../modules/stage-logs";
import { createdRandTime } from "./utils";

const createdAt = createdRandTime();

export const dummyLogBlock: LogBlock.AsObject = {
  index: 0,
  log: "This is stage log",
  severity: LogSeverity.SUCCESS,
  createdAt: createdAt.unix(),
};

export function createLogBlockFromObject(o: LogBlock.AsObject): LogBlock {
  const block = new LogBlock();
  block.setIndex(o.index);
  block.setLog(o.log);
  block.setSeverity(o.severity);
  block.setCreatedAt(o.createdAt);
  return block;
}
