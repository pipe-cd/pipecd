import { apiClient, apiRequest } from "./client";
import {
  GetMeRequest,
  GetMeResponse,
} from "pipe/pkg/app/web/api_client/service_pb";

export const getMe = (): Promise<GetMeResponse.AsObject> => {
  const req = new GetMeRequest();
  return apiRequest(req, apiClient.getMe);
};
