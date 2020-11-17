import { rest } from "msw";
import { GetDeploymentResponse } from "pipe/pkg/app/web/api_client/service_pb";
import {
  Commit,
  Deployment,
  DeploymentTrigger,
  PipelineStage,
} from "pipe/pkg/app/web/model/deployment_pb";
import {
  ApplicationGitPath,
  ApplicationGitRepository,
} from "pipe/pkg/app/web/model/common_pb";
import { dummyDeployment } from "../../__fixtures__/dummy-deployment";
import { serialize } from "../serializer";
import { createMask } from "../utils";

const createGitPath = (o: ApplicationGitPath.AsObject): ApplicationGitPath => {
  const gitPath = new ApplicationGitPath();
  gitPath.setConfigFilename(o.configFilename);
  gitPath.setConfigPath(o.configPath);
  gitPath.setPath(o.path);
  gitPath.setUrl(o.url);
  if (o.repo) {
    const repo = new ApplicationGitRepository();
    repo.setId(o.repo.id);
    repo.setBranch(o.repo.branch);
    repo.setRemote(o.repo.remote);
    gitPath.setRepo(repo);
  }
  return gitPath;
};

const createTrigger = (o: DeploymentTrigger.AsObject): DeploymentTrigger => {
  const trigger = new DeploymentTrigger();
  trigger.setCommander(o.commander);
  trigger.setTimestamp(o.timestamp);
  if (o.commit) {
    const commit = new Commit();
    commit.setAuthor(o.commit.author);
    commit.setBranch(o.commit.branch);
    commit.setCreatedAt(o.commit.createdAt);
    commit.setHash(o.commit.hash);
    commit.setMessage(o.commit.message);
    commit.setPullRequest(o.commit.pullRequest);
    commit.setUrl(o.commit.url);
    trigger.setCommit(commit);
  }
  return trigger;
};

const createPipeline = (o: PipelineStage.AsObject[]): PipelineStage[] => {
  return [];
};

const createDeployment = (o: Deployment.AsObject): Deployment => {
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
  o.gitPath && deployment.setGitPath(createGitPath(o.gitPath));
  o.trigger && deployment.setTrigger(createTrigger(o.trigger));
  o.stagesList && deployment.setStagesList(createPipeline(o.stagesList));
  return deployment;
};

export const deploymentHandlers = [
  rest.post<Uint8Array>(createMask("/GetDeployment"), (req, res, ctx) => {
    const response = new GetDeploymentResponse();

    response.setDeployment(createDeployment(dummyDeployment));

    const data = serialize(response.serializeBinary());
    return res(
      ctx.status(200),
      ctx.set("Content-Type", "application/grpc-web+proto"),
      ctx.body(data)
    );
  }),
];
