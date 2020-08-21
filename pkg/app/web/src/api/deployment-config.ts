import { apiClient, apiRequest } from "./client";
import {
  ListDeploymentConfigTemplatesRequest,
  ListDeploymentConfigTemplatesResponse,
} from "pipe/pkg/app/web/api_client/service_pb";

export const getDeploymentConfigTemplates = ({
  applicationId,
  labelsList,
}: ListDeploymentConfigTemplatesRequest.AsObject): Promise<
  ListDeploymentConfigTemplatesResponse.AsObject
> => {
  const req = new ListDeploymentConfigTemplatesRequest();
  req.setApplicationId(applicationId);
  req.setLabelsList(labelsList);
  return apiRequest(req, apiClient.listDeploymentConfigTemplates);
};
