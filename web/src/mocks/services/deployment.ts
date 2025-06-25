import {
  CancelDeploymentResponse,
  GetDeploymentResponse,
  ListDeploymentsResponse,
} from "pipecd/web/api_client/service_pb";
import {
  createDeploymentFromObject,
  dummyDeployment,
} from "~/__fixtures__/dummy-deployment";
import { createHandler } from "../create-handler";

export const deploymentHandlers = [
  createHandler<GetDeploymentResponse>("/GetDeployment", () => {
    const response = new GetDeploymentResponse();
    response.setDeployment(createDeploymentFromObject(dummyDeployment));
    return response;
  }),
  createHandler<ListDeploymentsResponse>("/ListDeployments", () => {
    const response = new ListDeploymentsResponse();
    response.setDeploymentsList([createDeploymentFromObject(dummyDeployment)]);
    return response;
  }),
  createHandler<CancelDeploymentResponse>("/CancelDeployment", () => {
    const response = new CancelDeploymentResponse();
    response.setCommandId("123");
    return response;
  }),
];
