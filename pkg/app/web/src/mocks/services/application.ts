import {
  AddApplicationResponse,
  DeleteApplicationResponse,
  DisableApplicationResponse,
  EnableApplicationResponse,
  GetApplicationResponse,
  ListApplicationsResponse,
  SyncApplicationResponse,
  UpdateApplicationResponse,
} from "pipe/pkg/app/web/api_client/service_pb";
import {
  createApplicationFromObject,
  dummyApplication,
} from "../../__fixtures__/dummy-application";
import { createHandler } from "../create-handler";

export const applicationHandlers = [
  createHandler<SyncApplicationResponse>("/SyncApplication", () => {
    const response = new SyncApplicationResponse();
    response.setCommandId("sync-command");
    return response;
  }),
  createHandler<EnableApplicationResponse>("/EnableApplication", () => {
    return new EnableApplicationResponse();
  }),
  createHandler<DisableApplicationResponse>("/DisableApplication", () => {
    return new DisableApplicationResponse();
  }),
  createHandler<DeleteApplicationResponse>("/DeleteApplication", () => {
    return new DeleteApplicationResponse();
  }),
  createHandler<AddApplicationResponse>("/AddApplication", () => {
    const response = new AddApplicationResponse();
    response.setApplicationId(dummyApplication.id);
    return response;
  }),
  createHandler<UpdateApplicationResponse>("/UpdateApplication", () => {
    return new UpdateApplicationResponse();
  }),
  createHandler<ListApplicationsResponse>("/ListApplications", () => {
    const response = new ListApplicationsResponse();
    response.setApplicationsList([
      createApplicationFromObject(dummyApplication),
    ]);
    return response;
  }),
  createHandler<GetApplicationResponse>("/GetApplication", () => {
    const response = new GetApplicationResponse();
    response.setApplication(createApplicationFromObject(dummyApplication));
    return response;
  }),
];
