import { dummyDeployment } from "./dummy-deployment";
import { createRandTimes, randomUUID } from "./utils";
import { ListDeploymentTracesResponse } from "~~/api_client/service_pb";

const [createdAt, completedAt] = createRandTimes(3);

export const dummyDeploymentTrace: ListDeploymentTracesResponse.DeploymentTraceRes.AsObject = {
  trace: {
    id: randomUUID(),
    title: "title",
    author: "user",
    commitTimestamp: createdAt.unix(),
    commitMessage: "commit-message",
    commitHash: "commit-hash",
    commitUrl: "commit-url",
    createdAt: createdAt.unix(),
    updatedAt: completedAt.unix(),
    completedAt: completedAt.unix(),
  },
  deploymentsList: [dummyDeployment],
};
