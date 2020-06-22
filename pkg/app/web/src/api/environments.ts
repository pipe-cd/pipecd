import { apiClient, apiRequest } from "./client";
import {
  ListEnvironmentsRequest,
  ListEnvironmentsResponse,
} from "pipe/pkg/app/web/api_client/service_pb";

export const getEnvironments = ({}: ListEnvironmentsRequest.AsObject): Promise<
  ListEnvironmentsResponse.AsObject
> => {
  const req = new ListEnvironmentsRequest();
  return apiRequest(req, apiClient.listEnvironments);
};
