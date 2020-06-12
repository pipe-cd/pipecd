import { apiClient, apiRequest } from "./client";
import {
  GetDeploymentRequest,
  GetDeploymentResponse,
  ListDeploymentsRequest,
  ListDeploymentsResponse,
} from "pipe/pkg/app/web/api_client/service_pb";

export const getDeployment = ({
  deploymentId,
}: GetDeploymentRequest.AsObject): Promise<GetDeploymentResponse.AsObject> => {
  const req = new GetDeploymentRequest();
  req.setDeploymentId(deploymentId);
  return apiRequest(req, apiClient.getDeployment);
};

export const getDeployments = ({}: ListDeploymentsRequest.AsObject): Promise<
  ListDeploymentsResponse.AsObject
> => {
  const req = new ListDeploymentsRequest();
  return apiRequest(req, apiClient.listDeployments);
};
