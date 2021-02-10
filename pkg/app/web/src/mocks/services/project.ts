import {
  UpdateProjectStaticAdminResponse,
  GetProjectResponse,
} from "pipe/pkg/app/web/api_client/service_pb";
import { Project } from "pipe/pkg/app/web/model/project_pb";
import { dummyProject } from "../../__fixtures__/dummy-project";
import { createHandler } from "../create-handler";

function createProjectFromObject(o: Project.AsObject): Project {
  const project = new Project();
  project.setId(o.id);
  project.setCreatedAt(o.createdAt);
  project.setUpdatedAt(o.updatedAt);
  return project;
}

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
