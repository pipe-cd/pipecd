import { ApplicationKind } from "pipe/pkg/app/web/model/common_pb";
import {
  Deployment,
  DeploymentStatus,
  PipelineStage,
} from "../modules/deployments";
import { createGitPathFromObject } from "./common";
import { dummyApplication } from "./dummy-application";
import { dummyEnv } from "./dummy-environment";
import { dummyPiped } from "./dummy-piped";
import { dummyStage } from "./dummy-stage";
import { dummyTrigger, createTriggerFromObject } from "./dummy-trigger";
import { createdRandTime, subtractRandTimeFrom } from "./utils";
import faker from "faker";

faker.seed(1);

const completedAt = createdRandTime();
const updatedAt = subtractRandTimeFrom(completedAt);
const createdAt = subtractRandTimeFrom(updatedAt);

export const dummyDeployment: Deployment.AsObject = {
  id: faker.random.uuid(),
  pipedId: dummyPiped.id,
  projectId: "project-1",
  applicationName: dummyApplication.name,
  applicationId: dummyApplication.id,
  runningCommitHash: faker.random.uuid().slice(0, 8),
  stagesList: [dummyStage],
  status: DeploymentStatus.DEPLOYMENT_SUCCESS,
  statusReason: "good",
  trigger: dummyTrigger,
  version: "0.0.0",
  cloudProvider: "kube-1",
  createdAt: createdAt.unix(),
  updatedAt: updatedAt.unix(),
  completedAt: completedAt.unix(),
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

function createPipelineFromObject(
  o: PipelineStage.AsObject[]
): PipelineStage[] {
  return [];
}

export function createDeploymentFromObject(o: Deployment.AsObject): Deployment {
  const deployment = new Deployment();
  deployment.setId(o.id);
  deployment.setApplicationId(o.applicationId);
  deployment.setApplicationName(o.applicationName);
  deployment.setCloudProvider(o.cloudProvider);
  deployment.setCompletedAt(o.completedAt);
  deployment.setCreatedAt(o.createdAt);
  deployment.setEnvId(o.envId);
  deployment.setKind(o.kind);
  deployment.setPipedId(o.pipedId);
  deployment.setProjectId(o.projectId);
  deployment.setRunningCommitHash(o.runningCommitHash);
  deployment.setStatus(o.status);
  deployment.setStatusReason(o.statusReason);
  deployment.setSummary(o.summary);
  deployment.setUpdatedAt(o.updatedAt);
  deployment.setVersion(o.version);
  o.gitPath && deployment.setGitPath(createGitPathFromObject(o.gitPath));
  o.trigger && deployment.setTrigger(createTriggerFromObject(o.trigger));
  o.stagesList &&
    deployment.setStagesList(createPipelineFromObject(o.stagesList));
  return deployment;
}
