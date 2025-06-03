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
import { InsightRange } from "./insight.config";
import { determineTimeRange } from "~/utils/determine-time-range";

export const useInsightDeploymentFrequency = (
  filterValues: {
    labels: string[];
    range: InsightRange;
    resolution: InsightResolution;
    applicationId: string;
  },
  queryOption: UseQueryOptions<InsightDataPoint.AsObject[]> = {}
): UseQueryResult<InsightDataPoint.AsObject[]> => {
  return useQuery({
    queryKey: ["insight", "deployment-frequency", filterValues],
    queryFn: async () => {
      const labels = new Array<[string, string]>();
      if (filterValues.labels) {
        for (const label of filterValues.labels) {
          const pair = label.split(":");
          if (pair.length === 2) labels.push([pair[0], pair[1]]);
        }
      }

      const [rangeFrom, rangeTo] = determineTimeRange(
        filterValues.range,
        filterValues.resolution
      );

      const data = await InsightAPI.getInsightData({
        metricsKind: InsightMetricsKind.DEPLOYMENT_FREQUENCY,
        rangeFrom: rangeFrom,
        rangeTo: rangeTo,
        resolution: filterValues.resolution,
        applicationId: filterValues.applicationId,
        labelsMap: labels,
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
