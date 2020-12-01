import { apiClient, apiRequest } from "./client";
import {
  GenerateAPIKeyRequest,
  GenerateAPIKeyResponse,
  DisableAPIKeyRequest,
  DisableAPIKeyResponse,
  ListAPIKeysRequest,
  ListAPIKeysResponse,
} from "pipe/pkg/app/web/api_client/service_pb";

export const getAPIKeys = ({
  options,
}: ListAPIKeysRequest.AsObject): Promise<ListAPIKeysResponse.AsObject> => {
  const req = new ListAPIKeysRequest();
  if (options) {
    const opt = new ListAPIKeysRequest.Options();
    opt.setEnabled(options.enabled);
    req.setOptions(opt);
  }
  return apiRequest(req, apiClient.listAPIKeys);
};

export const generateAPIKey = ({
  name,
  role,
}: GenerateAPIKeyRequest.AsObject): Promise<
  GenerateAPIKeyResponse.AsObject
> => {
  const req = new GenerateAPIKeyRequest();
  req.setName(name);
  req.setRole(role);
  return apiRequest(req, apiClient.generateAPIKey);
};

export const disableAPIKey = ({
  id,
}: DisableAPIKeyRequest.AsObject): Promise<DisableAPIKeyResponse.AsObject> => {
  const req = new DisableAPIKeyRequest();
  req.setId(id);
  return apiRequest(req, apiClient.disableAPIKey);
};
