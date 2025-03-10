import {
  ListDeploymentTracesRequest,
  ListDeploymentTracesResponse,
} from "~~/api_client/service_pb";
import { apiClient, apiRequest } from "./client";

export const getDeploymentTraces = ({
  options,
  pageSize,
  cursor,
  pageMinUpdatedAt,
}: ListDeploymentTracesRequest.AsObject): Promise<
  ListDeploymentTracesResponse.AsObject
> => {
  const req = new ListDeploymentTracesRequest();
  if (options) {
    const opts = new ListDeploymentTracesRequest.Options();
    opts.setCommitHash(options.commitHash);
    req.setOptions(opts);
    req.setPageSize(pageSize);
    req.setCursor(cursor);
    req.setPageMinUpdatedAt(pageMinUpdatedAt);
  }
  return apiRequest(req, apiClient.listDeploymentTraces);
};
