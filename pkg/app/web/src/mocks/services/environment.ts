import { StatusCode } from "grpc-web";
import {
  ListEnvironmentsResponse,
  DeleteEnvironmentResponse,
  AddEnvironmentResponse,
} from "pipe/pkg/app/web/api_client/service_pb";
import {
  createEnvFromObject,
  deletedDummyEnv,
  dummyEnv,
} from "~/__fixtures__/dummy-environment";
import { createHandler, createHandlerWithError } from "../create-handler";

export const listEnvironmentHandler = createHandler<ListEnvironmentsResponse>(
  "/ListEnvironments",
  () => {
    const response = new ListEnvironmentsResponse();
    response.setEnvironmentsList([
      createEnvFromObject(dummyEnv),
      createEnvFromObject(deletedDummyEnv),
    ]);
    return response;
  }
);

export const addEnvironmentHandler = createHandler<AddEnvironmentResponse>(
  "/AddEnvironment",
  () => {
    return new AddEnvironmentResponse();
  }
);

export const deleteEnvironmentHandler = createHandler<
  DeleteEnvironmentResponse
>("/DeleteEnvironment", () => {
  return new DeleteEnvironmentResponse();
});

export const deleteEnvironmentFailedHandler = createHandlerWithError(
  "/DeleteEnvironment",
  StatusCode.FAILED_PRECONDITION
);

export const environmentHandlers = [
  listEnvironmentHandler,
  deleteEnvironmentHandler,
  addEnvironmentHandler,
];
