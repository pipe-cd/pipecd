import {
  AddApplicationResponse,
  DeleteApplicationResponse,
  DisableApplicationResponse,
  EnableApplicationResponse,
  GetApplicationRequest,
  GetApplicationResponse,
  ListApplicationsResponse,
  SyncApplicationResponse,
  UpdateApplicationResponse,
} from "pipecd/web/api_client/service_pb";
import { ApplicationKind } from "~/modules/applications";
import {
  createApplicationFromObject,
  dummyApplication,
  dummyApps,
} from "~/__fixtures__/dummy-application";
import { createHandler } from "../create-handler";

export const updateApplicationHandler = createHandler<
  UpdateApplicationResponse
>("/UpdateApplication", () => {
  return new UpdateApplicationResponse();
});

export const listApplicationsHandler = createHandler<ListApplicationsResponse>(
  "/ListApplications",
  () => {
    const response = new ListApplicationsResponse();
    response.setApplicationsList([
      createApplicationFromObject(dummyApps[ApplicationKind.KUBERNETES]),
      createApplicationFromObject(dummyApps[ApplicationKind.TERRAFORM]),
      createApplicationFromObject(dummyApps[ApplicationKind.LAMBDA]),
      createApplicationFromObject(dummyApps[ApplicationKind.CLOUDRUN]),
      createApplicationFromObject(dummyApps[ApplicationKind.ECS]),
    ]);
    return response;
  }
);

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
  updateApplicationHandler,
  createHandler<AddApplicationResponse>("/AddApplication", () => {
    const response = new AddApplicationResponse();
    response.setApplicationId(dummyApplication.id);
    return response;
  }),
  createHandler<UpdateApplicationResponse>("/UpdateApplication", () => {
    return new UpdateApplicationResponse();
  }),
  listApplicationsHandler,
  createHandler<GetApplicationResponse>("/GetApplication", (requestBody) => {
    const response = new GetApplicationResponse();
    const params = GetApplicationRequest.deserializeBinary(requestBody);
    const appId = params.getApplicationId();
    const findApp = Object.values(dummyApps).find((app) => app.id === appId);

    response.setApplication(
      createApplicationFromObject(findApp ?? dummyApplication)
    );
    return response;
  }),
];
