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
  SyncApplicationRequest,
  SyncApplicationResponse,
  DisableApplicationRequest,
  DisableApplicationResponse,
  EnableApplicationRequest,
  EnableApplicationResponse,
  UpdateApplicationRequest,
  UpdateApplicationResponse,
  DeleteApplicationRequest,
  DeleteApplicationResponse,
  UpdateApplicationDescriptionRequest,
  UpdateApplicationDescriptionResponse,
  ListUnregisteredApplicationsRequest,
  ListUnregisteredApplicationsResponse,
} from "pipe/pkg/app/web/api_client/service_pb";
import { ApplicationGitPath } from "pipe/pkg/app/web/model/common_pb";
import { ApplicationGitRepository } from "pipe/pkg/app/web/model/common_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";

export const getApplicationLiveState = ({
  applicationId,
}: GetApplicationLiveStateRequest.AsObject): Promise<
  GetApplicationLiveStateResponse.AsObject
> => {
  const req = new GetApplicationLiveStateRequest();
  req.setApplicationId(applicationId);
  return apiRequest(req, apiClient.getApplicationLiveState);
};

export const getApplications = ({
  options,
}: ListApplicationsRequest.AsObject): Promise<
  ListApplicationsResponse.AsObject
> => {
  const req = new ListApplicationsRequest();
  if (options) {
    const o = new ListApplicationsRequest.Options();
    o.setEnvIdsList(options.envIdsList);
    o.setKindsList(options.kindsList);
    o.setSyncStatusesList(options.syncStatusesList);
    o.setName(options.name);
    if (options.enabled !== undefined) {
      const enabled = new google_protobuf_wrappers_pb.BoolValue();
      enabled.setValue((options.enabled.value as unknown) as boolean);
      o.setEnabled(enabled);
    }
    req.setOptions(o);
  }
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
  const repository = new ApplicationGitRepository();
  if (gitPath.repo) {
    repository.setId(gitPath.repo.id);
    repository.setBranch(gitPath.repo.branch);
    repository.setRemote(gitPath.repo.remote);
    appGitPath.setRepo(repository);
  }
  appGitPath.setPath(gitPath.path);
  if (gitPath.configFilename && gitPath.configFilename !== "") {
    appGitPath.setConfigFilename(gitPath.configFilename);
  }
  req.setGitPath(appGitPath);
  return apiRequest(req, apiClient.addApplication);
};

export const syncApplication = async ({
  applicationId,
  syncStrategy,
}: SyncApplicationRequest.AsObject): Promise<
  SyncApplicationResponse.AsObject
> => {
  const req = new SyncApplicationRequest();
  req.setApplicationId(applicationId);
  req.setSyncStrategy(syncStrategy);
  return apiRequest(req, apiClient.syncApplication);
};

export const disableApplication = async ({
  applicationId,
}: DisableApplicationRequest.AsObject): Promise<
  DisableApplicationResponse.AsObject
> => {
  const req = new DisableApplicationRequest();
  req.setApplicationId(applicationId);
  return apiRequest(req, apiClient.disableApplication);
};

export const enableApplication = async ({
  applicationId,
}: EnableApplicationRequest.AsObject): Promise<
  EnableApplicationResponse.AsObject
> => {
  const req = new EnableApplicationRequest();
  req.setApplicationId(applicationId);
  return apiRequest(req, apiClient.enableApplication);
};

export const updateApplication = async ({
  applicationId,
  cloudProvider,
  envId,
  kind,
  name,
  pipedId,
}: Required<UpdateApplicationRequest.AsObject>): Promise<
  UpdateApplicationResponse.AsObject
> => {
  const req = new UpdateApplicationRequest();
  req.setApplicationId(applicationId);
  req.setName(name);
  req.setEnvId(envId);
  req.setPipedId(pipedId);
  req.setCloudProvider(cloudProvider);
  req.setKind(kind);
  return apiRequest(req, apiClient.updateApplication);
};

export const deleteApplication = async ({
  applicationId,
}: DeleteApplicationRequest.AsObject): Promise<
  DeleteApplicationResponse.AsObject
> => {
  const req = new DeleteApplicationRequest();
  req.setApplicationId(applicationId);
  return apiRequest(req, apiClient.deleteApplication);
};

export const updateDescription = async ({
  applicationId,
  description,
}: UpdateApplicationDescriptionRequest.AsObject): Promise<
  UpdateApplicationDescriptionResponse.AsObject
> => {
  const req = new UpdateApplicationDescriptionRequest();
  req.setApplicationId(applicationId);
  req.setDescription(description);
  return apiRequest(req, apiClient.updateApplicationDescription);
};

export const getUnregisteredApplications = ({
}: ListUnregisteredApplicationsRequest.AsObject): Promise<
  ListUnregisteredApplicationsResponse.AsObject
> => {
  const req = new ListUnregisteredApplicationsRequest();
  return apiRequest(req, apiClient.listUnregisteredApplications);
};
