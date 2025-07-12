import {
  useInfiniteQuery,
  UseInfiniteQueryResult,
} from "@tanstack/react-query";
import {
  ListDeploymentTracesRequest,
  ListDeploymentTracesResponse,
} from "~~/api_client/service_pb";
import * as deploymentTracesApi from "~/api/deploymentTraces";

// 30 days
const TIME_RANGE_LIMIT_IN_SECONDS = 2592000;
const ITEMS_PER_PAGE = 50;
const FETCH_MORE_ITEMS_PER_PAGE = 30;

export type DeploymentTraceFilterOptions = {
  commitHash?: string;
};

const convertFilterOptions = (
  options: DeploymentTraceFilterOptions
): ListDeploymentTracesRequest.Options.AsObject => {
  return {
    commitHash: options?.commitHash || "",
  };
};

export const useGetDeploymentTracesInfinite = (
  options: DeploymentTraceFilterOptions
): UseInfiniteQueryResult<{
  tracesList: ListDeploymentTracesResponse.DeploymentTraceRes.AsObject[];
  cursor: string;
  minUpdatedAt: number;
  hasMore: boolean;
}> => {
  return useInfiniteQuery({
    queryKey: ["deployment-traces", options],
    queryFn: async ({ pageParam }) => {
      const isFirstPage = !pageParam;
      const minUpdatedAt =
        pageParam?.minUpdatedAt ??
        Math.round(Date.now() / 1000 - TIME_RANGE_LIMIT_IN_SECONDS);

      const pageSize = isFirstPage ? ITEMS_PER_PAGE : FETCH_MORE_ITEMS_PER_PAGE;

      const {
        tracesList,
        cursor,
      } = await deploymentTracesApi.getDeploymentTraces({
        options: convertFilterOptions(options),
        pageSize: pageSize,
        cursor: pageParam?.cursor || "",
        pageMinUpdatedAt: minUpdatedAt,
      });

      return {
        tracesList,
        cursor: cursor || pageParam?.cursor || "",
        minUpdatedAt,
        hasMore: tracesList.length >= pageSize,
      };
    },
    getNextPageParam: (lastPage, allPages) => {
      const isFirstPage = allPages.length === 0;

      const ITEMS_PER_PAGE_TO_USE = isFirstPage
        ? ITEMS_PER_PAGE
        : FETCH_MORE_ITEMS_PER_PAGE;

      const isHasMore = lastPage.tracesList.length >= ITEMS_PER_PAGE_TO_USE;

      const initMinUpdateAt = Math.round(
        Date.now() / 1000 - TIME_RANGE_LIMIT_IN_SECONDS
      );
      let minUpdatedAt = initMinUpdateAt;

      if (isHasMore) {
        minUpdatedAt = lastPage.minUpdatedAt || initMinUpdateAt;
      } else {
        minUpdatedAt = lastPage.minUpdatedAt - TIME_RANGE_LIMIT_IN_SECONDS;
      }

      return {
        hasMore: isHasMore,
        cursor: lastPage.cursor,
        minUpdatedAt: minUpdatedAt,
      };
    },
    refetchOnWindowFocus: false,
  });
};
