import {
  UpdateProjectStaticAdminResponse,
  GetProjectResponse,
} from "pipe/pkg/app/web/api_client/service_pb";
import { Project } from "pipe/pkg/app/web/model/project_pb";
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
    response.setProject(new Project());
    return response;
  }),
];
