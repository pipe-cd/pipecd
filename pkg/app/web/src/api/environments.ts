import { apiClient, apiRequest } from "./client";
import {
  ListEnvironmentsRequest,
  ListEnvironmentsResponse,
  AddEnvironmentRequest,
  AddEnvironmentResponse,
} from "pipe/pkg/app/web/api_client/service_pb";

export const getEnvironments = (): Promise<
  ListEnvironmentsResponse.AsObject
> => {
  const req = new ListEnvironmentsRequest();
  return apiRequest(req, apiClient.listEnvironments);
};

export const AddEnvironment = ({
  name,
  desc,
}: AddEnvironmentRequest.AsObject): Promise<
  AddEnvironmentResponse.AsObject
> => {
  const req = new AddEnvironmentRequest();
  req.setName(name);
  req.setDesc(desc);
  return apiRequest(req, apiClient.addEnvironment);
};
