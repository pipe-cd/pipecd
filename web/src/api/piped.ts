import { apiClient, apiRequest } from "./client";
import {
  RegisterPipedRequest,
  RegisterPipedResponse,
  ListPipedsRequest,
  ListPipedsResponse,
  DisablePipedRequest,
  DisablePipedResponse,
  EnablePipedRequest,
  EnablePipedResponse,
  RecreatePipedKeyRequest,
  RecreatePipedKeyResponse,
  GenerateApplicationSealedSecretRequest,
  GenerateApplicationSealedSecretResponse,
  UpdatePipedRequest,
  UpdatePipedResponse,
  DeleteOldPipedKeysRequest,
  DeleteOldPipedKeysResponse,
  UpdatePipedDesiredVersionRequest,
  UpdatePipedDesiredVersionResponse,
  ListReleasedVersionsResponse,
  ListReleasedVersionsRequest,
  RestartPipedRequest,
  RestartPipedResponse,
  ListDeprecatedNotesRequest,
  ListDeprecatedNotesResponse,
} from "pipecd/web/api_client/service_pb";

export const getPipeds = ({
  withStatus,
}: ListPipedsRequest.AsObject): Promise<ListPipedsResponse.AsObject> => {
  const req = new ListPipedsRequest();
  req.setWithStatus(withStatus);
  return apiRequest(req, apiClient.listPipeds);
};

export const listReleasedVersions = (): Promise<
  ListReleasedVersionsResponse.AsObject
> => {
  const req = new ListReleasedVersionsRequest();
  return apiRequest(req, apiClient.listReleasedVersions);
};

export const listBreakingChanges = ({
  projectId,
}: ListDeprecatedNotesRequest.AsObject): Promise<
  ListDeprecatedNotesResponse.AsObject
> => {
  const req = new ListDeprecatedNotesRequest();
  req.setProjectId(projectId);
  return apiRequest(req, apiClient.listDeprecatedNotes);
};

export const registerPiped = ({
  name,
  desc,
}: RegisterPipedRequest.AsObject): Promise<RegisterPipedResponse.AsObject> => {
  const req = new RegisterPipedRequest();
  req.setName(name);
  req.setDesc(desc);
  return apiRequest(req, apiClient.registerPiped);
};

export const disablePiped = ({
  pipedId,
}: DisablePipedRequest.AsObject): Promise<DisablePipedResponse.AsObject> => {
  const req = new DisablePipedRequest();
  req.setPipedId(pipedId);
  return apiRequest(req, apiClient.disablePiped);
};

export const enablePiped = ({
  pipedId,
}: EnablePipedRequest.AsObject): Promise<EnablePipedResponse.AsObject> => {
  const req = new EnablePipedRequest();
  req.setPipedId(pipedId);
  return apiRequest(req, apiClient.enablePiped);
};

export const restartPiped = ({
  pipedId,
}: RestartPipedRequest.AsObject): Promise<RestartPipedResponse.AsObject> => {
  const req = new RestartPipedRequest();
  req.setPipedId(pipedId);
  return apiRequest(req, apiClient.restartPiped);
};

export const recreatePipedKey = ({
  id,
}: RecreatePipedKeyRequest.AsObject): Promise<
  RecreatePipedKeyResponse.AsObject
> => {
  const req = new RecreatePipedKeyRequest();
  req.setId(id);
  return apiRequest(req, apiClient.recreatePipedKey);
};

export const deleteOldPipedKey = ({
  pipedId,
}: DeleteOldPipedKeysRequest.AsObject): Promise<
  DeleteOldPipedKeysResponse.AsObject
> => {
  const req = new DeleteOldPipedKeysRequest();
  req.setPipedId(pipedId);
  return apiRequest(req, apiClient.deleteOldPipedKeys);
};

export const generateApplicationSealedSecret = ({
  pipedId,
  data,
  base64Encoding,
}: GenerateApplicationSealedSecretRequest.AsObject): Promise<
  GenerateApplicationSealedSecretResponse.AsObject
> => {
  const req = new GenerateApplicationSealedSecretRequest();
  req.setPipedId(pipedId);
  req.setData(data);
  req.setBase64Encoding(base64Encoding);
  return apiRequest(req, apiClient.generateApplicationSealedSecret);
};

export const updatePiped = ({
  pipedId,
  name,
  desc,
}: UpdatePipedRequest.AsObject): Promise<UpdatePipedResponse.AsObject> => {
  const req = new UpdatePipedRequest();
  req.setPipedId(pipedId);
  req.setName(name);
  req.setDesc(desc);
  return apiRequest(req, apiClient.updatePiped);
};

export const updatePipedDesiredVersion = ({
  version,
  pipedIdsList,
}: UpdatePipedDesiredVersionRequest.AsObject): Promise<
  UpdatePipedDesiredVersionResponse.AsObject
> => {
  const req = new UpdatePipedDesiredVersionRequest();
  req.setVersion(version);
  req.setPipedIdsList(pipedIdsList);
  return apiRequest(req, apiClient.updatePipedDesiredVersion);
};
