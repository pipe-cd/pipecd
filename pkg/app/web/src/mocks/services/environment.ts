import { StatusCode } from "grpc-web";
import {
  ListEnvironmentsResponse,
  DeleteEnvironmentResponse,
} from "pipe/pkg/app/web/api_client/service_pb";
import {
  createEnvFromObject,
  dummyEnv,
} from "../../__fixtures__/dummy-environment";
import { createHandler, createHandlerWithError } from "../create-handler";

export const deleteEnvironmentHandler = createHandler<
  DeleteEnvironmentResponse
>("/DeleteEnvironment", () => {
  const response = new DeleteEnvironmentResponse();
  return response;
});

export const deleteEnvironmentFailedHandler = createHandlerWithError(
  "/DeleteEnvironment",
  StatusCode.FAILED_PRECONDITION
);

export const listEnvironmentHandler = createHandler<ListEnvironmentsResponse>(
  "/ListEnvironments",
  () => {
    const response = new ListEnvironmentsResponse();
    response.setEnvironmentsList([createEnvFromObject(dummyEnv)]);
    return response;
  }
);

export const environmentHandlers = [
  listEnvironmentHandler,
  deleteEnvironmentHandler,
];
