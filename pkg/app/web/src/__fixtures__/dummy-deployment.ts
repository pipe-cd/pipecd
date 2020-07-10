import { Deployment, DeploymentStatus } from "../modules/deployments";
import { ApplicationKind } from "pipe/pkg/app/web/model/common_pb";
import { dummyApplication } from "./dummy-application";
import { dummyEnv } from "./dummy-environment";

export const dummyDeployment: Deployment = {
  id: "deployment-1",
  pipedId: "piped-1",
  projectId: "project-1",
  applicationName: dummyApplication.name,
  applicationId: dummyApplication.id,
  runningCommitHash: "123456abcdefg",
  stagesList: [],
  status: DeploymentStatus.DEPLOYMENT_SUCCESS,
  statusDescription: "good",
  trigger: {
    commander: "user",
    timestamp: 1,
    commit: {
      author: "user",
      branch: "branch",
      createdAt: 1,
      hash: "12345abc",
      message: "fix",
      pullRequest: 123,
    },
  },
  updatedAt: 1,
  version: "0.0.0",
  cloudProvider: "kube-1",
  completedAt: 1,
  createdAt: 1,
  description: "description",
  envId: dummyEnv.id,
  gitPath: { configPath: "", path: "", repoId: "" },
  kind: ApplicationKind.KUBERNETES,
  metadataMap: [],
};
