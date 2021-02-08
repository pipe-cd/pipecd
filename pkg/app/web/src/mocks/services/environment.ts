import { ListEnvironmentsResponse } from "pipe/pkg/app/web/api_client/service_pb";
import {
  createEnvFromObject,
  dummyEnv,
} from "../../__fixtures__/dummy-environment";
import { createHandler } from "../create-handler";

export const environmentHandlers = [
  createHandler<ListEnvironmentsResponse>("/ListEnvironments", () => {
    const response = new ListEnvironmentsResponse();
    response.setEnvironmentsList([createEnvFromObject(dummyEnv)]);
    return response;
  }),
];
