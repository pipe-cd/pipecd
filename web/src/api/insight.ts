import { apiClient, apiRequest } from "./client";
import {
  GetInsightDataRequest,
  GetInsightDataResponse,
  GetInsightApplicationCountRequest,
  GetInsightApplicationCountResponse,
} from "pipecd/web/api_client/service_pb";

export const getInsightData = ({
  applicationId,
  metricsKind,
  rangeFrom,
  rangeTo,
}: GetInsightDataRequest.AsObject): Promise<
  GetInsightDataResponse.AsObject
> => {
  const req = new GetInsightDataRequest();
  req.setApplicationId(applicationId);
  req.setMetricsKind(metricsKind);
  req.setRangeFrom(rangeFrom);
  req.setRangeTo(rangeTo);
  return apiRequest(req, apiClient.getInsightData);
};

export const getApplicationCount = (): Promise<
  GetInsightApplicationCountResponse.AsObject
> => {
  const req = new GetInsightApplicationCountRequest();
  return apiRequest(req, apiClient.getInsightApplicationCount);
};
