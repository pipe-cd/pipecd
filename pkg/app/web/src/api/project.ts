import { apiClient, apiRequest } from "./client";
import {
  GetProjectRequest,
  GetProjectResponse,
  UpdateProjectStaticAdminRequest,
  UpdateProjectStaticAdminResponse,
  UpdateProjectSSOConfigRequest,
  UpdateProjectSSOConfigResponse,
  EnableStaticAdminRequest,
  EnableStaticAdminResponse,
  DisableStaticAdminRequest,
  DisableStaticAdminResponse,
} from "pipe/pkg/app/web/api_client/service_pb";
import { ProjectSSOConfig } from "pipe/pkg/app/web/model/project_pb";

export const getProject = (): Promise<GetProjectResponse.AsObject> => {
  const req = new GetProjectRequest();
  return apiRequest(req, apiClient.getProject);
};

export const updateStaticAdminPassword = ({
  password,
}: {
  password: string;
}): Promise<UpdateProjectStaticAdminResponse.AsObject> => {
  const req = new UpdateProjectStaticAdminRequest();
  req.setPassword(password);
  return apiRequest(req, apiClient.updateProjectStaticAdmin);
};

export const updateStaticAdminUsername = ({
  username,
}: {
  username: string;
}): Promise<UpdateProjectStaticAdminResponse.AsObject> => {
  const req = new UpdateProjectStaticAdminRequest();
  req.setUsername(username);
  return apiRequest(req, apiClient.updateProjectStaticAdmin);
};

export const enableStaticAdmin = (): Promise<
  EnableStaticAdminResponse.AsObject
> => {
  const req = new EnableStaticAdminRequest();
  return apiRequest(req, apiClient.enableStaticAdmin);
};

export const disableStaticAdmin = (): Promise<
  DisableStaticAdminResponse.AsObject
> => {
  const req = new DisableStaticAdminRequest();
  return apiRequest(req, apiClient.disableStaticAdmin);
};

export const updateGitHubSSO = ({
  clientId,
  clientSecret,
  baseUrl,
  uploadUrl,
}: ProjectSSOConfig.GitHub.AsObject): Promise<
  UpdateProjectSSOConfigResponse.AsObject
> => {
  const req = new UpdateProjectSSOConfigRequest();
  const params = new ProjectSSOConfig();
  const github = new ProjectSSOConfig.GitHub();
  github.setClientId(clientId);
  github.setClientSecret(clientSecret);
  github.setBaseUrl(baseUrl);
  github.setUploadUrl(uploadUrl);
  req.setSso(params);
  return apiRequest(req, apiClient.updateProjectSSOConfig);
};
