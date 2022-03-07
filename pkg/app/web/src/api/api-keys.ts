import { apiClient, apiRequest } from "./client";
import {
  GenerateAPIKeyRequest,
  GenerateAPIKeyResponse,
  DisableAPIKeyRequest,
  DisableAPIKeyResponse,
  ListAPIKeysRequest,
  ListAPIKeysResponse,
} from "pipecd/pkg/app/web/api_client/service_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";

export const getAPIKeys = ({
  options,
}: {
  options: {
    enabled: boolean;
  };
}): Promise<ListAPIKeysResponse.AsObject> => {
  const req = new ListAPIKeysRequest();
  const opt = new ListAPIKeysRequest.Options();
  const enabled = new google_protobuf_wrappers_pb.BoolValue();
  enabled.setValue(options.enabled);
  opt.setEnabled(enabled);
  req.setOptions(opt);
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
