import {
  DeploymentTrigger,
  SyncStrategy,
  Commit,
} from "pipe/pkg/app/web/model/deployment_pb";
import { createdRandTime } from "./utils";
import faker from "faker";

const commitTimestamp = createdRandTime();

export const dummyTrigger: DeploymentTrigger.AsObject = {
  commander: "user",
  timestamp: commitTimestamp.unix(),
  commit: {
    author: "pipecd-user",
    branch: "feat/awesome-feature",
    createdAt: 1,
    hash: faker.random.uuid().slice(0, 8),
    message: "fix",
    pullRequest: 123,
    url: "",
  },
  syncStrategy: SyncStrategy.AUTO,
};

function createCommitFromObject(o: Commit.AsObject): Commit {
  const commit = new Commit();
  commit.setAuthor(o.author);
  commit.setBranch(o.branch);
  commit.setCreatedAt(o.createdAt);
  commit.setHash(o.hash);
  commit.setMessage(o.message);
  commit.setPullRequest(o.pullRequest);
  commit.setUrl(o.url);
  return commit;
}

export function createTriggerFromObject(
  o: DeploymentTrigger.AsObject
): DeploymentTrigger {
  const trigger = new DeploymentTrigger();
  trigger.setCommander(o.commander);
  trigger.setTimestamp(o.timestamp);
  trigger.setSyncStrategy(o.syncStrategy);
  if (o.commit) {
    trigger.setCommit(createCommitFromObject(o.commit));
  }
  return trigger;
}
