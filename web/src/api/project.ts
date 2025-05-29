import {
  DisableStaticAdminRequest,
  DisableStaticAdminResponse,
  EnableStaticAdminRequest,
  EnableStaticAdminResponse,
  GetProjectRequest,
  GetProjectResponse,
  UpdateProjectRBACConfigRequest,
  UpdateProjectRBACConfigResponse,
  UpdateProjectSSOConfigRequest,
  UpdateProjectSSOConfigResponse,
  UpdateProjectStaticAdminRequest,
  UpdateProjectStaticAdminResponse,
  AddProjectUserGroupRequest,
  AddProjectUserGroupResponse,
  DeleteProjectUserGroupRequest,
  DeleteProjectUserGroupResponse,
  AddProjectRBACRoleRequest,
  AddProjectRBACRoleResponse,
  DeleteProjectRBACRoleRequest,
  DeleteProjectRBACRoleResponse,
  UpdateProjectRBACRoleRequest,
  UpdateProjectRBACRoleResponse,
} from "pipecd/web/api_client/service_pb";
import {
  ProjectRBACConfig,
  ProjectSSOConfig,
  ProjectRBACPolicy,
} from "pipecd/web/model/project_pb";
import { apiClient, apiRequest } from "./client";

export const getProject = (): Promise<GetProjectResponse.AsObject> => {
  const req = new GetProjectRequest();
  return apiRequest(req, apiClient.getProject);
};

export const updateStaticAdmin = ({
  username,
  password,
}: {
  username?: string;
  password?: string;
}): Promise<UpdateProjectStaticAdminResponse.AsObject> => {
  const req = new UpdateProjectStaticAdminRequest();
  if (username) req.setUsername(username);
  if (password) req.setPassword(password);
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

export const updateRBAC = ({
  admin,
  editor,
  viewer,
}: ProjectRBACConfig.AsObject): Promise<
  UpdateProjectRBACConfigResponse.AsObject
> => {
  const req = new UpdateProjectRBACConfigRequest();
  const rbac = new ProjectRBACConfig();
  rbac.setAdmin(admin);
  rbac.setEditor(editor);
  rbac.setViewer(viewer);
  req.setRbac(rbac);
  return apiRequest(req, apiClient.updateProjectRBACConfig);
};

export const updateGitHubSSO = ({
  clientId,
  clientSecret,
  baseUrl,
  uploadUrl,
}: {
  clientId: string;
  clientSecret: string;
  baseUrl?: string;
  uploadUrl?: string;
}): Promise<UpdateProjectSSOConfigResponse.AsObject> => {
  const req = new UpdateProjectSSOConfigRequest();
  const sso = new ProjectSSOConfig();
  const github = new ProjectSSOConfig.GitHub();
  github.setClientId(clientId);
  github.setClientSecret(clientSecret);
  if (baseUrl) github.setBaseUrl(baseUrl);
  if (uploadUrl) github.setUploadUrl(uploadUrl);

  sso.setGithub(github);
  req.setSso(sso);
  return apiRequest(req, apiClient.updateProjectSSOConfig);
};

export const addUserGroup = ({
  ssoGroup,
  role,
}: AddProjectUserGroupRequest.AsObject): Promise<
  AddProjectUserGroupResponse.AsObject
> => {
  const req = new AddProjectUserGroupRequest();
  req.setSsoGroup(ssoGroup);
  req.setRole(role);
  return apiRequest(req, apiClient.addProjectUserGroup);
};

export const deleteUserGroup = ({
  ssoGroup,
}: DeleteProjectUserGroupRequest.AsObject): Promise<
  DeleteProjectUserGroupResponse.AsObject
> => {
  const req = new DeleteProjectUserGroupRequest();
  req.setSsoGroup(ssoGroup);
  return apiRequest(req, apiClient.deleteProjectUserGroup);
};

export const addRBACRole = ({
  name,
  policies,
}: {
  name: string;
  policies: ProjectRBACPolicy[];
}): Promise<AddProjectRBACRoleResponse.AsObject> => {
  const req = new AddProjectRBACRoleRequest();
  req.setName(name);
  req.setPoliciesList(policies);
  return apiRequest(req, apiClient.addProjectRBACRole);
};

export const deleteRBACRole = ({
  name,
}: DeleteProjectRBACRoleRequest.AsObject): Promise<
  DeleteProjectRBACRoleResponse.AsObject
> => {
  const req = new DeleteProjectRBACRoleRequest();
  req.setName(name);
  return apiRequest(req, apiClient.deleteProjectRBACRole);
};

export const updateRBACRole = ({
  name,
  policies,
}: {
  name: string;
  policies: ProjectRBACPolicy[];
}): Promise<UpdateProjectRBACRoleResponse.AsObject> => {
  const req = new UpdateProjectRBACRoleRequest();
  req.setName(name);
  req.setPoliciesList(policies);
  return apiRequest(req, apiClient.updateProjectRBACRole);
};
