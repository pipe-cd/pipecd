import { apiClient, apiRequest } from "./client";
import {
  RegisterPipedRequest,
  RegisterPipedResponse,
} from "pipe/pkg/app/web/api_client/service_pb";

export const registerPiped = ({
  desc,
}: RegisterPipedRequest.AsObject): Promise<RegisterPipedResponse.AsObject> => {
  const req = new RegisterPipedRequest();
  req.setDesc(desc);
  return apiRequest(req, apiClient.registerPiped);
};
