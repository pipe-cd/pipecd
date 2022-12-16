import { apiClient, apiRequest } from "./client";
import {
  GetInsightDataRequest,
  GetInsightDataResponse,
  GetInsightApplicationCountRequest,
  GetInsightApplicationCountResponse,
} from "pipecd/web/api_client/service_pb";
import { InsightResultType } from "pipecd/web/model/insight_pb";

export const getInsightData = ({
  metricsKind,
  rangeFrom,
  rangeTo,
  resolution,
  applicationId,
  labelsMap,
}: GetInsightDataRequest.AsObject): Promise<
  GetInsightDataResponse.AsObject
> => {
  const req = new GetInsightDataRequest();
  req.setMetricsKind(metricsKind);

  // Convert unix milli second to unix second.
  rangeFrom = Math.floor(rangeFrom / 1000);
  rangeTo = Math.floor(rangeTo / 1000);
  req.setRangeFrom(rangeFrom);
  req.setRangeTo(rangeTo);
  req.setResolution(resolution);

  req.setResolution(resolution);
  req.setApplicationId(applicationId);
  for (const label of labelsMap) {
    req.getLabelsMap().set(label[0], label[1]);
  }

  const p = apiRequest(req, apiClient.getInsightData) as Promise<
    GetInsightDataResponse.AsObject
  >;
  return p.then((value) => {
    // Server uses unix second and Client uses unix milli second.
    // So we convert it here and above.
    // But it might be better to unify client and server.
    return convertTimestamp(value);
  });
};

// Convert unix sec to unix milli
const convertTimestamp = (
  value: GetInsightDataResponse.AsObject
): GetInsightDataResponse.AsObject => {
  switch (value.type) {
    case InsightResultType.MATRIX:
      return {
        ...value,
        updatedAt: value.updatedAt * 1000,
        matrixList: value.matrixList.map((sampleStream) => ({
          labelsMap: sampleStream.labelsMap,
          dataPointsList: sampleStream.dataPointsList.map((value) => ({
            timestamp: value.timestamp * 1000,
            value: value.value,
          })),
        })),
      };
    case InsightResultType.VECTOR:
      return {
        ...value,
        updatedAt: value.updatedAt * 1000,
        vectorList: value.vectorList.map((insightSample) => ({
          ...insightSample,
          dataPoint:
            insightSample.dataPoint !== undefined
              ? {
                  timestamp: insightSample.dataPoint.timestamp * 1000,
                  value: insightSample.dataPoint.value,
                }
              : undefined,
        })),
      };
  }
};

export const getApplicationCount = (): Promise<
  GetInsightApplicationCountResponse.AsObject
> => {
  const req = new GetInsightApplicationCountRequest();
  return apiRequest(req, apiClient.getInsightApplicationCount);
};
