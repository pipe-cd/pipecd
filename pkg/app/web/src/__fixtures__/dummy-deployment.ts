import { Deployment, DeploymentStatus } from "../modules/deployments";
import { ApplicationKind } from "pipe/pkg/app/web/model/common_pb";
import { dummyApplication } from "./dummy-application";
import { dummyEnv } from "./dummy-environment";
import { dummyPiped } from "./dummy-piped";
import { dummyStage } from "./dummy-stage";

export const dummyDeployment: Deployment = {
  id: "deployment-1",
  pipedId: dummyPiped.id,
  projectId: "project-1",
  applicationName: dummyApplication.name,
  applicationId: dummyApplication.id,
  runningCommitHash: "123456abcdefg",
  stagesList: [dummyStage],
  status: DeploymentStatus.DEPLOYMENT_SUCCESS,
  statusReason: "good",
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
      url: "",
    },
  },
  updatedAt: 1,
  version: "0.0.0",
  cloudProvider: "kube-1",
  completedAt: 1,
  createdAt: 1,
  summary:
    "Quick sync by deploying the new version and configuring all traffic to it because no pipeline was configured",
  envId: dummyEnv.id,
  gitPath: {
    configPath: "",
    configFilename: "",
    path: "",
    url: "",
    repo: {
      id: "repo-1",
      branch: "master",
      remote: "xxx",
    },
  },
  kind: ApplicationKind.KUBERNETES,
  metadataMap: [],
};
