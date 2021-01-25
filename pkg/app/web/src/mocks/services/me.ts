import { GetMeResponse } from "pipe/pkg/app/web/api_client/service_pb";
import { Role } from "../../modules/me";
import { createHandler } from "../create-handler";

export const meHandlers = [
  createHandler<GetMeResponse>("/GetMe", () => {
    const response = new GetMeResponse();
    response.setAvatarUrl("https://test.pipecd.dev/avatar.jpg");
    response.setSubject("hello-pipecd");
    response.setProjectId("pipecd");
    response.setProjectRole(Role.ProjectRole.ADMIN);
    return response;
  }),
];
