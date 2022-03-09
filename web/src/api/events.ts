import { apiClient, apiRequest } from "./client";
import {
  ListEventsRequest,
  ListEventsResponse,
} from "pipecd/web/api_client/service_pb";

export const getEvents = ({
  options,
  pageSize,
  cursor,
  pageMinUpdatedAt,
}: ListEventsRequest.AsObject): Promise<ListEventsResponse.AsObject> => {
  const req = new ListEventsRequest();
  if (options) {
    const opts = new ListEventsRequest.Options();
    opts.setStatusesList(options.statusesList);
    opts.setName(options.name);
    options.labelsMap.forEach((label) => {
      opts.getLabelsMap().set(label[0], label[1]);
    });
    req.setOptions(opts);
    req.setPageSize(pageSize);
    req.setCursor(cursor);
    req.setPageMinUpdatedAt(pageMinUpdatedAt);
  }
  return apiRequest(req, apiClient.listEvents);
};
