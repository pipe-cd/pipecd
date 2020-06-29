import { apiClient, apiRequest } from "./client";
import {
  RegisterPipedRequest,
  RegisterPipedResponse,
  ListPipedsRequest,
  ListPipedsResponse,
} from "pipe/pkg/app/web/api_client/service_pb";

export const getPipeds = ({
  withStatus,
}: ListPipedsRequest.AsObject): Promise<ListPipedsResponse.AsObject> => {
  const req = new ListPipedsRequest();
  req.setWithStatus(withStatus);
  return apiRequest(req, apiClient.listPipeds);
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
