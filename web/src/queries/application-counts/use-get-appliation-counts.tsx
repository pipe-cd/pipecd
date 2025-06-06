import {
  useQuery,
  UseQueryOptions,
  UseQueryResult,
} from "@tanstack/react-query";
import { StatusCode } from "grpc-web";
import {
  ApplicationActiveStatus,
  ApplicationKind,
} from "pipecd/web/model/common_pb";
import { InsightApplicationCountLabelKey } from "pipecd/web/model/insight_pb";
import { getApplicationCount } from "~/api/insight";
import { APPLICATION_ACTIVE_STATUS_NAME } from "~/constants/application-active-status";
import { APPLICATION_KIND_TEXT } from "~/constants/application-kind";

export const INSIGHT_APPLICATION_COUNT_LABEL_KEY_TEXT: Record<
  InsightApplicationCountLabelKey,
  string
> = {
  [InsightApplicationCountLabelKey.KIND]: "KIND",
  [InsightApplicationCountLabelKey.ACTIVE_STATUS]: "ACTIVE_STATUS",
};

export interface ApplicationCounts {
  updatedAt: number;
  counts: Record<string, Record<string, number>>;
  summary: {
    total: number;
    enabled: number;
    disabled: number;
  };
}

const createInitialCount = (): Record<string, number> => ({
  [APPLICATION_ACTIVE_STATUS_NAME[ApplicationActiveStatus.ENABLED]]: 0,
  [APPLICATION_ACTIVE_STATUS_NAME[ApplicationActiveStatus.DISABLED]]: 0,
});

const createInitialCounts = (): Record<string, Record<string, number>> => ({
  [APPLICATION_KIND_TEXT[ApplicationKind.KUBERNETES]]: createInitialCount(),
  [APPLICATION_KIND_TEXT[ApplicationKind.TERRAFORM]]: createInitialCount(),
  [APPLICATION_KIND_TEXT[ApplicationKind.LAMBDA]]: createInitialCount(),
  [APPLICATION_KIND_TEXT[ApplicationKind.CLOUDRUN]]: createInitialCount(),
  [APPLICATION_KIND_TEXT[ApplicationKind.ECS]]: createInitialCount(),
});

const initialState: ApplicationCounts = {
  updatedAt: 0,
  counts: createInitialCounts(),
  summary: { total: 0, enabled: 0, disabled: 0 },
};

export const useGetApplicationCounts = (
  queryOption: UseQueryOptions<ApplicationCounts> = {}
): UseQueryResult<ApplicationCounts> => {
  return useQuery({
    queryKey: ["applications", "count"],
    queryFn: async () => {
      const res = await getApplicationCount().catch((e: { code: number }) => {
        // NOT_FOUND is the initial normal state, so it is excluded from error handling.
        if (e.code !== StatusCode.NOT_FOUND) {
          throw e;
        }
      });

      if (!res) {
        // Handling of errors with code NOT_FOUND.
        return initialState;
      }

      const counts: Record<
        string,
        Record<string, number>
      > = createInitialCounts();

      res.countsList.forEach((count) => {
        const [, kindName] =
          count.labelsMap.find(
            (val) =>
              val[0] ===
              INSIGHT_APPLICATION_COUNT_LABEL_KEY_TEXT[
                InsightApplicationCountLabelKey.KIND
              ]
          ) || [];
        const [, activeStatusName] =
          count.labelsMap.find(
            (val) =>
              val[0] ===
              INSIGHT_APPLICATION_COUNT_LABEL_KEY_TEXT[
                InsightApplicationCountLabelKey.ACTIVE_STATUS
              ]
          ) || [];
        if (!kindName || !activeStatusName) {
          return;
        }

        counts[kindName][activeStatusName] = count.count;
      });

      const summary = {
        total: 0,
        enabled: 0,
        disabled: 0,
      };

      res.countsList.forEach((count) => {
        summary.total += count.count;
        const [, activeStatusName] =
          count.labelsMap.find(
            (val) =>
              val[0] ===
              INSIGHT_APPLICATION_COUNT_LABEL_KEY_TEXT[
                InsightApplicationCountLabelKey.ACTIVE_STATUS
              ]
          ) || [];
        if (!activeStatusName) return;

        if (
          activeStatusName ===
          APPLICATION_ACTIVE_STATUS_NAME[ApplicationActiveStatus.ENABLED]
        ) {
          summary.enabled += count.count;
        }
        if (
          activeStatusName ===
          APPLICATION_ACTIVE_STATUS_NAME[ApplicationActiveStatus.DISABLED]
        ) {
          summary.disabled += count.count;
        }
      });

      return { updatedAt: res.updatedAt, counts, summary };
    },
    ...queryOption,
  });
};
