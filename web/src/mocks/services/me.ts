import { GetMeResponse } from "pipecd/web/api_client/service_pb";
import { dummyMe } from "~/__fixtures__/dummy-me";
import { createHandler } from "../create-handler";

export const getMeHandler = createHandler<GetMeResponse>("/GetMe", () => {
  const response = new GetMeResponse();
  response.setSubject(dummyMe.subject);
  response.setProjectId(dummyMe.projectId);
  response.setAvatarUrl(dummyMe.avatarUrl);
  response.setProjectRole(dummyMe.projectRole);
  return response;
});

export const meHandlers = [getMeHandler];
