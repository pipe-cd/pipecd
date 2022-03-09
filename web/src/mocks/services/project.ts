import {
  GetProjectResponse,
  UpdateProjectStaticAdminResponse,
} from "pipecd/web/api_client/service_pb";
import {
  createProjectFromObject,
  dummyProject,
} from "~/__fixtures__/dummy-project";
import { createHandler } from "../create-handler";

export const projectHandlers = [
  createHandler<UpdateProjectStaticAdminResponse>(
    "/UpdateProjectStaticAdmin",
    () => {
      return new UpdateProjectStaticAdminResponse();
    }
  ),
  createHandler<GetProjectResponse>("/GetProject", () => {
    const response = new GetProjectResponse();
    response.setProject(createProjectFromObject(dummyProject));
    return response;
  }),
];
