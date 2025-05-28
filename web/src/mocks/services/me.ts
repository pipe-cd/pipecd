import { GetMeResponse } from "pipecd/web/api_client/service_pb";
import { dummyMe } from "~/__fixtures__/dummy-me";
import { createHandler, createHandlerWithError } from "../create-handler";
import { StatusCode } from "grpc-web";

export const getMeHandler = createHandler<GetMeResponse>("/GetMe", () => {
  const response = new GetMeResponse();
  response.setSubject(dummyMe.subject);
  response.setProjectId(dummyMe.projectId);
  response.setAvatarUrl(dummyMe.avatarUrl);
  return response;
});

export const getMeUnauthenticatedHandler = createHandlerWithError(
  "/GetMe",
  StatusCode.UNAUTHENTICATED
);

export const meHandlers = [getMeHandler];
