import { Deployment, DeploymentStatus } from "../modules/deployments";
import { ApplicationKind } from "pipe/pkg/app/web/model/common_pb";

export const dummyDeployment: Deployment = {
  id: "deployment-1",
  pipedId: "piped-1",
  projectId: "project-1",
  runningCommitHash: "123456abcdefg",
  stagesList: [],
  status: DeploymentStatus.DEPLOYMENT_SUCCESS,
  statusDescription: "good",
  trigger: {
    commander: "user",
    timestamp: 0,
    commit: {
      author: "user",
      branch: "branch",
      createdAt: 0,
      hash: "12345abc",
      message: "fix",
      pullRequest: 123,
    },
  },
  updatedAt: 0,
  version: "0.0.0",
  applicationId: "app-1",
  cloudProvider: "kube-1",
  completedAt: 0,
  createdAt: 0,
  description: "description",
  envId: "env-1",
  gitPath: { configPath: "", path: "", repoId: "" },
  kind: ApplicationKind.KUBERNETES,
  metadataMap: [],
};
