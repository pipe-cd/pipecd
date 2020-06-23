import { apiClient, apiRequest } from "./client";
import {
  GetApplicationLiveStateRequest,
  GetApplicationLiveStateResponse,
  GetApplicationRequest,
  GetApplicationResponse,
  ListApplicationsRequest,
  ListApplicationsResponse,
  AddApplicationRequest,
  AddApplicationResponse,
} from "pipe/pkg/app/web/api_client/service_pb";
import { ApplicationGitPath } from "pipe/pkg/app/web/model/common_pb";

export const getApplicationLiveState = ({
  applicationId,
}: GetApplicationLiveStateRequest.AsObject): Promise<
  GetApplicationLiveStateResponse.AsObject
> => {
  const req = new GetApplicationLiveStateRequest();
  req.setApplicationId(applicationId);
  return apiRequest(req, apiClient.getApplicationLiveState);
};

export const getApplications = (): Promise<
  ListApplicationsResponse.AsObject
> => {
  const req = new ListApplicationsRequest();
  return apiRequest(req, apiClient.listApplications);
};

export const getApplication = ({
  applicationId,
}: GetApplicationRequest.AsObject): Promise<
  GetApplicationResponse.AsObject
> => {
  const req = new GetApplicationRequest();
  req.setApplicationId(applicationId);
  return apiRequest(req, apiClient.getApplication);
};

export const addApplication = async ({
  name,
  envId,
  pipedId,
  cloudProvider,
  kind,
  gitPath,
}: Required<AddApplicationRequest.AsObject>): Promise<
  AddApplicationResponse.AsObject
> => {
  const req = new AddApplicationRequest();
  req.setName(name);
  req.setEnvId(envId);
  req.setPipedId(pipedId);
  req.setCloudProvider(cloudProvider);
  req.setKind(kind);
  const appGitPath = new ApplicationGitPath();
  appGitPath.setRepoId(gitPath.repoId);
  appGitPath.setPath(gitPath.path);
  if (gitPath.configPath && gitPath.configPath !== "") {
    appGitPath.setConfigPath(gitPath.configPath);
  }
  req.setGitPath(appGitPath);
  return apiRequest(req, apiClient.addApplication);
};
