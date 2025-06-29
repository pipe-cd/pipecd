import {
  useQuery,
  UseQueryOptions,
  UseQueryResult,
} from "@tanstack/react-query";
import * as InsightAPI from "~/api/insight";
import {
  InsightDataPoint,
  InsightMetricsKind,
  InsightResolution,
  InsightResultType,
} from "pipecd/web/model/insight_pb";
import dayjs from "dayjs";

export const useInsightDeploymentFrequency24h = (
  queryOption: UseQueryOptions<InsightDataPoint.AsObject[]> = {}
): UseQueryResult<InsightDataPoint.AsObject[]> => {
  return useQuery({
    queryKey: ["insight", "deployment-frequency-24h"],
    queryFn: async () => {
      const rangeTo = dayjs.utc().endOf("day").valueOf();
      const rangeFrom = dayjs.utc().startOf("day").valueOf();

      const data = await InsightAPI.getInsightData({
        metricsKind: InsightMetricsKind.DEPLOYMENT_FREQUENCY,
        rangeFrom,
        rangeTo,
        resolution: InsightResolution.DAILY,
        applicationId: "",
        labelsMap: [],
      });
      if (data.type == InsightResultType.MATRIX) {
        return data.matrixList[0].dataPointsList;
      } else {
        return [];
      }
    },
    ...queryOption,
  });
};
