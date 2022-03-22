import { apiClient, apiRequest } from "./client";
import { GetMeRequest, GetMeResponse } from "pipecd/web/api_client/service_pb";

export const getMe = (): Promise<GetMeResponse.AsObject> => {
  const req = new GetMeRequest();
  return apiRequest(req, apiClient.getMe);
};
