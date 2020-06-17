import { apiClient, apiRequest } from "./client";
import {
  GetApplicationLiveStateRequest,
  GetApplicationLiveStateResponse,
  ListApplicationsRequest,
  ListApplicationsResponse,
} from "pipe/pkg/app/web/api_client/service_pb";

export const getApplicationLiveState = ({
  applicationId,
}: GetApplicationLiveStateRequest.AsObject): Promise<
  GetApplicationLiveStateResponse.AsObject
> => {
  const req = new GetApplicationLiveStateRequest();
  req.setApplicationId(applicationId);
  return apiRequest(req, apiClient.getApplicationLiveState);
};

export const getApplications = (): Promise<
  ListApplicationsResponse.AsObject
> => {
  const req = new ListApplicationsRequest();
  return apiRequest(req, apiClient.listApplications);
};
