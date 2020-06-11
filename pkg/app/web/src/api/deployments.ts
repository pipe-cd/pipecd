import { apiClient, apiRequest } from "./client";
import {
  GetDeploymentRequest,
  GetDeploymentResponse,
} from "pipe/pkg/app/web/api_client/service_pb";

export const getDeployment = ({
  deploymentId,
}: GetDeploymentRequest.AsObject): Promise<GetDeploymentResponse.AsObject> => {
  const req = new GetDeploymentRequest();
  req.setDeploymentId(deploymentId);
  return apiRequest(req, apiClient.getDeployment);
};
