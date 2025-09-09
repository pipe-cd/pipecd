import { LogBlock, LogSeverity } from "~~/model/logblock_pb";
import { createRandTimes, randomWords } from "./utils";

const logTimes = createRandTimes(3);

export const dummyLogBlock: LogBlock.AsObject = {
  index: 0,
  log: randomWords(8),
  severity: LogSeverity.SUCCESS,
  createdAt: logTimes[0].unix(),
};

export const dummyLogBlocks: LogBlock.AsObject[] = [
  dummyLogBlock,
  {
    index: 1,
    log: randomWords(8),
    severity: LogSeverity.INFO,
    createdAt: logTimes[1].unix(),
  },
  {
    index: 2,
    log: randomWords(8),
    severity: LogSeverity.ERROR,
    createdAt: logTimes[2].unix(),
  },
];

export function createLogBlockFromObject(o: LogBlock.AsObject): LogBlock {
  const block = new LogBlock();
  block.setIndex(o.index);
  block.setLog(o.log);
  block.setSeverity(o.severity);
  block.setCreatedAt(o.createdAt);
  return block;
}
