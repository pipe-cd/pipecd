import { apiClient, apiRequest } from "./client";
import {
  GetInsightDataRequest,
  GetInsightDataResponse,
} from "pipe/pkg/app/web/api_client/service_pb";

export const getInsightData = ({
  applicationId,
  dataPointCount,
  metricsKind,
  rangeFrom,
  step,
}: GetInsightDataRequest.AsObject): Promise<
  GetInsightDataResponse.AsObject
> => {
  const req = new GetInsightDataRequest();
  req.setApplicationId(applicationId);
  req.setDataPointCount(dataPointCount);
  req.setMetricsKind(metricsKind);
  req.setRangeFrom(rangeFrom);
  req.setStep(step);
  return apiRequest(req, apiClient.getInsightData);
};
