import * as grpcWeb from 'grpc-web';

import * as pkg_app_server_service_webservice_service_pb from 'pipecd/web/app/server/service/webservice/service_pb';


export class WebServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  registerPiped(
    request: pkg_app_server_service_webservice_service_pb.RegisterPipedRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.RegisterPipedResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.RegisterPipedResponse>;

  updatePiped(
    request: pkg_app_server_service_webservice_service_pb.UpdatePipedRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.UpdatePipedResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.UpdatePipedResponse>;

  recreatePipedKey(
    request: pkg_app_server_service_webservice_service_pb.RecreatePipedKeyRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.RecreatePipedKeyResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.RecreatePipedKeyResponse>;

  deleteOldPipedKeys(
    request: pkg_app_server_service_webservice_service_pb.DeleteOldPipedKeysRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.DeleteOldPipedKeysResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.DeleteOldPipedKeysResponse>;

  enablePiped(
    request: pkg_app_server_service_webservice_service_pb.EnablePipedRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.EnablePipedResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.EnablePipedResponse>;

  disablePiped(
    request: pkg_app_server_service_webservice_service_pb.DisablePipedRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.DisablePipedResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.DisablePipedResponse>;

  listPipeds(
    request: pkg_app_server_service_webservice_service_pb.ListPipedsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.ListPipedsResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.ListPipedsResponse>;

  getPiped(
    request: pkg_app_server_service_webservice_service_pb.GetPipedRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.GetPipedResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.GetPipedResponse>;

  updatePipedDesiredVersion(
    request: pkg_app_server_service_webservice_service_pb.UpdatePipedDesiredVersionRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.UpdatePipedDesiredVersionResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.UpdatePipedDesiredVersionResponse>;

  restartPiped(
    request: pkg_app_server_service_webservice_service_pb.RestartPipedRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.RestartPipedResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.RestartPipedResponse>;

  listReleasedVersions(
    request: pkg_app_server_service_webservice_service_pb.ListReleasedVersionsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.ListReleasedVersionsResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.ListReleasedVersionsResponse>;

  listDeprecatedNotes(
    request: pkg_app_server_service_webservice_service_pb.ListDeprecatedNotesRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.ListDeprecatedNotesResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.ListDeprecatedNotesResponse>;

  addApplication(
    request: pkg_app_server_service_webservice_service_pb.AddApplicationRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.AddApplicationResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.AddApplicationResponse>;

  updateApplication(
    request: pkg_app_server_service_webservice_service_pb.UpdateApplicationRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.UpdateApplicationResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.UpdateApplicationResponse>;

  enableApplication(
    request: pkg_app_server_service_webservice_service_pb.EnableApplicationRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.EnableApplicationResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.EnableApplicationResponse>;

  disableApplication(
    request: pkg_app_server_service_webservice_service_pb.DisableApplicationRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.DisableApplicationResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.DisableApplicationResponse>;

  deleteApplication(
    request: pkg_app_server_service_webservice_service_pb.DeleteApplicationRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.DeleteApplicationResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.DeleteApplicationResponse>;

  listApplications(
    request: pkg_app_server_service_webservice_service_pb.ListApplicationsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.ListApplicationsResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.ListApplicationsResponse>;

  syncApplication(
    request: pkg_app_server_service_webservice_service_pb.SyncApplicationRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.SyncApplicationResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.SyncApplicationResponse>;

  getApplication(
    request: pkg_app_server_service_webservice_service_pb.GetApplicationRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.GetApplicationResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.GetApplicationResponse>;

  generateApplicationSealedSecret(
    request: pkg_app_server_service_webservice_service_pb.GenerateApplicationSealedSecretRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.GenerateApplicationSealedSecretResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.GenerateApplicationSealedSecretResponse>;

  listUnregisteredApplications(
    request: pkg_app_server_service_webservice_service_pb.ListUnregisteredApplicationsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.ListUnregisteredApplicationsResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.ListUnregisteredApplicationsResponse>;

  listDeployments(
    request: pkg_app_server_service_webservice_service_pb.ListDeploymentsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.ListDeploymentsResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.ListDeploymentsResponse>;

  getDeployment(
    request: pkg_app_server_service_webservice_service_pb.GetDeploymentRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.GetDeploymentResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.GetDeploymentResponse>;

  getStageLog(
    request: pkg_app_server_service_webservice_service_pb.GetStageLogRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.GetStageLogResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.GetStageLogResponse>;

  cancelDeployment(
    request: pkg_app_server_service_webservice_service_pb.CancelDeploymentRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.CancelDeploymentResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.CancelDeploymentResponse>;

  skipStage(
    request: pkg_app_server_service_webservice_service_pb.SkipStageRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.SkipStageResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.SkipStageResponse>;

  approveStage(
    request: pkg_app_server_service_webservice_service_pb.ApproveStageRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.ApproveStageResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.ApproveStageResponse>;

  getApplicationLiveState(
    request: pkg_app_server_service_webservice_service_pb.GetApplicationLiveStateRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.GetApplicationLiveStateResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.GetApplicationLiveStateResponse>;

  getProject(
    request: pkg_app_server_service_webservice_service_pb.GetProjectRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.GetProjectResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.GetProjectResponse>;

  updateProjectStaticAdmin(
    request: pkg_app_server_service_webservice_service_pb.UpdateProjectStaticAdminRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.UpdateProjectStaticAdminResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.UpdateProjectStaticAdminResponse>;

  enableStaticAdmin(
    request: pkg_app_server_service_webservice_service_pb.EnableStaticAdminRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.EnableStaticAdminResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.EnableStaticAdminResponse>;

  disableStaticAdmin(
    request: pkg_app_server_service_webservice_service_pb.DisableStaticAdminRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.DisableStaticAdminResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.DisableStaticAdminResponse>;

  updateProjectSSOConfig(
    request: pkg_app_server_service_webservice_service_pb.UpdateProjectSSOConfigRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.UpdateProjectSSOConfigResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.UpdateProjectSSOConfigResponse>;

  updateProjectRBACConfig(
    request: pkg_app_server_service_webservice_service_pb.UpdateProjectRBACConfigRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.UpdateProjectRBACConfigResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.UpdateProjectRBACConfigResponse>;

  getMe(
    request: pkg_app_server_service_webservice_service_pb.GetMeRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.GetMeResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.GetMeResponse>;

  addProjectRBACRole(
    request: pkg_app_server_service_webservice_service_pb.AddProjectRBACRoleRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.AddProjectRBACRoleResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.AddProjectRBACRoleResponse>;

  updateProjectRBACRole(
    request: pkg_app_server_service_webservice_service_pb.UpdateProjectRBACRoleRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.UpdateProjectRBACRoleResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.UpdateProjectRBACRoleResponse>;

  deleteProjectRBACRole(
    request: pkg_app_server_service_webservice_service_pb.DeleteProjectRBACRoleRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.DeleteProjectRBACRoleResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.DeleteProjectRBACRoleResponse>;

  addProjectUserGroup(
    request: pkg_app_server_service_webservice_service_pb.AddProjectUserGroupRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.AddProjectUserGroupResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.AddProjectUserGroupResponse>;

  deleteProjectUserGroup(
    request: pkg_app_server_service_webservice_service_pb.DeleteProjectUserGroupRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.DeleteProjectUserGroupResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.DeleteProjectUserGroupResponse>;

  getCommand(
    request: pkg_app_server_service_webservice_service_pb.GetCommandRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.GetCommandResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.GetCommandResponse>;

  generateAPIKey(
    request: pkg_app_server_service_webservice_service_pb.GenerateAPIKeyRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.GenerateAPIKeyResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.GenerateAPIKeyResponse>;

  disableAPIKey(
    request: pkg_app_server_service_webservice_service_pb.DisableAPIKeyRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.DisableAPIKeyResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.DisableAPIKeyResponse>;

  listAPIKeys(
    request: pkg_app_server_service_webservice_service_pb.ListAPIKeysRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.ListAPIKeysResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.ListAPIKeysResponse>;

  getInsightData(
    request: pkg_app_server_service_webservice_service_pb.GetInsightDataRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.GetInsightDataResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.GetInsightDataResponse>;

  getInsightApplicationCount(
    request: pkg_app_server_service_webservice_service_pb.GetInsightApplicationCountRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.GetInsightApplicationCountResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.GetInsightApplicationCountResponse>;

  listDeploymentChains(
    request: pkg_app_server_service_webservice_service_pb.ListDeploymentChainsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.ListDeploymentChainsResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.ListDeploymentChainsResponse>;

  getDeploymentChain(
    request: pkg_app_server_service_webservice_service_pb.GetDeploymentChainRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.GetDeploymentChainResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.GetDeploymentChainResponse>;

  listEvents(
    request: pkg_app_server_service_webservice_service_pb.ListEventsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: pkg_app_server_service_webservice_service_pb.ListEventsResponse) => void
  ): grpcWeb.ClientReadableStream<pkg_app_server_service_webservice_service_pb.ListEventsResponse>;

}

export class WebServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  registerPiped(
    request: pkg_app_server_service_webservice_service_pb.RegisterPipedRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.RegisterPipedResponse>;

  updatePiped(
    request: pkg_app_server_service_webservice_service_pb.UpdatePipedRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.UpdatePipedResponse>;

  recreatePipedKey(
    request: pkg_app_server_service_webservice_service_pb.RecreatePipedKeyRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.RecreatePipedKeyResponse>;

  deleteOldPipedKeys(
    request: pkg_app_server_service_webservice_service_pb.DeleteOldPipedKeysRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.DeleteOldPipedKeysResponse>;

  enablePiped(
    request: pkg_app_server_service_webservice_service_pb.EnablePipedRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.EnablePipedResponse>;

  disablePiped(
    request: pkg_app_server_service_webservice_service_pb.DisablePipedRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.DisablePipedResponse>;

  listPipeds(
    request: pkg_app_server_service_webservice_service_pb.ListPipedsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.ListPipedsResponse>;

  getPiped(
    request: pkg_app_server_service_webservice_service_pb.GetPipedRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.GetPipedResponse>;

  updatePipedDesiredVersion(
    request: pkg_app_server_service_webservice_service_pb.UpdatePipedDesiredVersionRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.UpdatePipedDesiredVersionResponse>;

  restartPiped(
    request: pkg_app_server_service_webservice_service_pb.RestartPipedRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.RestartPipedResponse>;

  listReleasedVersions(
    request: pkg_app_server_service_webservice_service_pb.ListReleasedVersionsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.ListReleasedVersionsResponse>;

  listDeprecatedNotes(
    request: pkg_app_server_service_webservice_service_pb.ListDeprecatedNotesRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.ListDeprecatedNotesResponse>;

  addApplication(
    request: pkg_app_server_service_webservice_service_pb.AddApplicationRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.AddApplicationResponse>;

  updateApplication(
    request: pkg_app_server_service_webservice_service_pb.UpdateApplicationRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.UpdateApplicationResponse>;

  enableApplication(
    request: pkg_app_server_service_webservice_service_pb.EnableApplicationRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.EnableApplicationResponse>;

  disableApplication(
    request: pkg_app_server_service_webservice_service_pb.DisableApplicationRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.DisableApplicationResponse>;

  deleteApplication(
    request: pkg_app_server_service_webservice_service_pb.DeleteApplicationRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.DeleteApplicationResponse>;

  listApplications(
    request: pkg_app_server_service_webservice_service_pb.ListApplicationsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.ListApplicationsResponse>;

  syncApplication(
    request: pkg_app_server_service_webservice_service_pb.SyncApplicationRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.SyncApplicationResponse>;

  getApplication(
    request: pkg_app_server_service_webservice_service_pb.GetApplicationRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.GetApplicationResponse>;

  generateApplicationSealedSecret(
    request: pkg_app_server_service_webservice_service_pb.GenerateApplicationSealedSecretRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.GenerateApplicationSealedSecretResponse>;

  listUnregisteredApplications(
    request: pkg_app_server_service_webservice_service_pb.ListUnregisteredApplicationsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.ListUnregisteredApplicationsResponse>;

  listDeployments(
    request: pkg_app_server_service_webservice_service_pb.ListDeploymentsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.ListDeploymentsResponse>;

  getDeployment(
    request: pkg_app_server_service_webservice_service_pb.GetDeploymentRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.GetDeploymentResponse>;

  getStageLog(
    request: pkg_app_server_service_webservice_service_pb.GetStageLogRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.GetStageLogResponse>;

  cancelDeployment(
    request: pkg_app_server_service_webservice_service_pb.CancelDeploymentRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.CancelDeploymentResponse>;

  skipStage(
    request: pkg_app_server_service_webservice_service_pb.SkipStageRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.SkipStageResponse>;

  approveStage(
    request: pkg_app_server_service_webservice_service_pb.ApproveStageRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.ApproveStageResponse>;

  getApplicationLiveState(
    request: pkg_app_server_service_webservice_service_pb.GetApplicationLiveStateRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.GetApplicationLiveStateResponse>;

  getProject(
    request: pkg_app_server_service_webservice_service_pb.GetProjectRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.GetProjectResponse>;

  updateProjectStaticAdmin(
    request: pkg_app_server_service_webservice_service_pb.UpdateProjectStaticAdminRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.UpdateProjectStaticAdminResponse>;

  enableStaticAdmin(
    request: pkg_app_server_service_webservice_service_pb.EnableStaticAdminRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.EnableStaticAdminResponse>;

  disableStaticAdmin(
    request: pkg_app_server_service_webservice_service_pb.DisableStaticAdminRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.DisableStaticAdminResponse>;

  updateProjectSSOConfig(
    request: pkg_app_server_service_webservice_service_pb.UpdateProjectSSOConfigRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.UpdateProjectSSOConfigResponse>;

  updateProjectRBACConfig(
    request: pkg_app_server_service_webservice_service_pb.UpdateProjectRBACConfigRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.UpdateProjectRBACConfigResponse>;

  getMe(
    request: pkg_app_server_service_webservice_service_pb.GetMeRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.GetMeResponse>;

  addProjectRBACRole(
    request: pkg_app_server_service_webservice_service_pb.AddProjectRBACRoleRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.AddProjectRBACRoleResponse>;

  updateProjectRBACRole(
    request: pkg_app_server_service_webservice_service_pb.UpdateProjectRBACRoleRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.UpdateProjectRBACRoleResponse>;

  deleteProjectRBACRole(
    request: pkg_app_server_service_webservice_service_pb.DeleteProjectRBACRoleRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.DeleteProjectRBACRoleResponse>;

  addProjectUserGroup(
    request: pkg_app_server_service_webservice_service_pb.AddProjectUserGroupRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.AddProjectUserGroupResponse>;

  deleteProjectUserGroup(
    request: pkg_app_server_service_webservice_service_pb.DeleteProjectUserGroupRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.DeleteProjectUserGroupResponse>;

  getCommand(
    request: pkg_app_server_service_webservice_service_pb.GetCommandRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.GetCommandResponse>;

  generateAPIKey(
    request: pkg_app_server_service_webservice_service_pb.GenerateAPIKeyRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.GenerateAPIKeyResponse>;

  disableAPIKey(
    request: pkg_app_server_service_webservice_service_pb.DisableAPIKeyRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.DisableAPIKeyResponse>;

  listAPIKeys(
    request: pkg_app_server_service_webservice_service_pb.ListAPIKeysRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.ListAPIKeysResponse>;

  getInsightData(
    request: pkg_app_server_service_webservice_service_pb.GetInsightDataRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.GetInsightDataResponse>;

  getInsightApplicationCount(
    request: pkg_app_server_service_webservice_service_pb.GetInsightApplicationCountRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.GetInsightApplicationCountResponse>;

  listDeploymentChains(
    request: pkg_app_server_service_webservice_service_pb.ListDeploymentChainsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.ListDeploymentChainsResponse>;

  getDeploymentChain(
    request: pkg_app_server_service_webservice_service_pb.GetDeploymentChainRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.GetDeploymentChainResponse>;

  listEvents(
    request: pkg_app_server_service_webservice_service_pb.ListEventsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<pkg_app_server_service_webservice_service_pb.ListEventsResponse>;

}

