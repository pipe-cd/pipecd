import { ApplicationKind } from "pipecd/web/model/common_pb";
import { Deployment, DeploymentStatus } from "~/modules/deployments";
import { createGitPathFromObject } from "./common";
import { dummyApplication } from "./dummy-application";
import { dummyPiped } from "./dummy-piped";
import { createPipelineFromObject, dummyPipeline } from "./dummy-pipeline";
import { createTriggerFromObject, dummyTrigger } from "./dummy-trigger";
import { createRandTimes, randomUUID } from "./utils";

const [createdAt, completedAt] = createRandTimes(3);

export const dummyDeployment: Deployment.AsObject = {
  id: randomUUID(),
  pipedId: dummyPiped.id,
  projectId: "project-1",
  applicationName: dummyApplication.name,
  applicationId: dummyApplication.id,
  runningCommitHash: randomUUID().slice(0, 8),
  runningConfigFilename: ".pipe.yaml",
  stagesList: dummyPipeline,
  status: DeploymentStatus.DEPLOYMENT_SUCCESS,
  statusReason: "good",
  trigger: dummyTrigger,
  version: "0.0.0",
  versionsList: [],
  cloudProvider: "kube-1",
  platformProvider: "kube-1",
  deployTargetsByPluginMap: [],
  labelsMap: [],
  createdAt: createdAt.unix(),
  updatedAt: completedAt.unix(),
  completedAt: completedAt.unix(),
  summary:
    "Quick sync by deploying the new version and configuring all traffic to it because no pipeline was configured",
  gitPath: {
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
  deploymentChainId: "",
  deploymentChainBlockIndex: 0,
};

export function createDeploymentFromObject(o: Deployment.AsObject): Deployment {
  const deployment = new Deployment();
  deployment.setId(o.id);
  deployment.setApplicationId(o.applicationId);
  deployment.setApplicationName(o.applicationName);
  deployment.setCloudProvider(o.cloudProvider);
  deployment.setCompletedAt(o.completedAt);
  deployment.setCreatedAt(o.createdAt);
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
