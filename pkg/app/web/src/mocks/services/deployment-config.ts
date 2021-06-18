import { ListDeploymentConfigTemplatesResponse } from "pipe/pkg/app/web/api_client/service_pb";
import {
  dummyDeploymentConfigTemplates,
  deploymentConfigTemplateFromObject,
} from "~/__fixtures__/dummy-deployment-config";
import { createHandler } from "../create-handler";

export const listDeploymentConfigTemplatesHandler = createHandler<
  ListDeploymentConfigTemplatesResponse
>("/ListDeploymentConfigTemplates", () => {
  const response = new ListDeploymentConfigTemplatesResponse();
  response.setTemplatesList(
    dummyDeploymentConfigTemplates.map(deploymentConfigTemplateFromObject)
  );
  return response;
});

export const deploymentConfigTemplatesHandlers = [
  listDeploymentConfigTemplatesHandler,
];
