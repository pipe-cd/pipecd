import { apiClient, apiRequest } from "./client";
import {
  ListEnvironmentsRequest,
  ListEnvironmentsResponse,
  UpdateEnvironmentDescRequest,
  UpdateEnvironmentDescResponse,
  DeleteEnvironmentRequest,
  DeleteEnvironmentResponse,
} from "pipe/pkg/app/web/api_client/service_pb";

export const getEnvironments = (): Promise<
  ListEnvironmentsResponse.AsObject
> => {
  const req = new ListEnvironmentsRequest();
  return apiRequest(req, apiClient.listEnvironments);
};

export const updateEnvironmentDesc = (
  desc: string
): Promise<UpdateEnvironmentDescResponse.AsObject> => {
  const req = new UpdateEnvironmentDescRequest();
  return apiRequest(req, apiClient.updateEnvironmentDesc);
};

export const deleteEnvironment = ({
  environmentId,
}: DeleteEnvironmentRequest.AsObject): DeleteEnvironmentResponse.AsObject => {
  const req = new DeleteEnvironmentRequest();
  req.setEnvironmentId(environmentId);
  return apiRequest(req, apiClient.deleteEnvironment);
};
