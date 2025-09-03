import {
  useInfiniteQuery,
  UseInfiniteQueryResult,
} from "@tanstack/react-query";
import {
  ListDeploymentTracesRequest,
  ListDeploymentTracesResponse,
} from "~~/api_client/service_pb";
import * as deploymentTracesApi from "~/api/deploymentTraces";
import { useCallback, useState } from "react";

// 30 days
const TIME_RANGE_LIMIT_IN_SECONDS = 2592000;
const FIRST_PAGE_SIZE = 50;
const FETCH_MORE_PAGE_SIZE = 30;

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

type QueryType = UseInfiniteQueryResult<{
  tracesList: ListDeploymentTracesResponse.DeploymentTraceRes.AsObject[];
  cursor: string;
  minUpdatedAt: number;
  hasMore: boolean;
}>;

export const useGetDeploymentTracesInfinite = (
  options: DeploymentTraceFilterOptions
): {
  data: QueryType["data"];
  isFetching: QueryType["isFetching"];
  fetchNextPage: () => void;
  refetch: () => void;
  isSuccess: boolean;
} => {
  const [localMinUpdatedAt, setLocalMinUpdatedAt] = useState(
    Math.round(Date.now() / 1000 - TIME_RANGE_LIMIT_IN_SECONDS)
  );

  const { fetchNextPage, data, ...otherParams } = useInfiniteQuery({
    queryKey: ["deployment-traces", options],
    queryFn: async ({
      pageParam,
    }: {
      pageParam?: {
        cursor?: string;
        minUpdatedAt?: number;
      };
    }) => {
      const isFirstPage = !pageParam;
      const minUpdatedAt = pageParam?.minUpdatedAt ?? localMinUpdatedAt;

      const pageSize = isFirstPage ? FIRST_PAGE_SIZE : FETCH_MORE_PAGE_SIZE;

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
    refetchOnWindowFocus: false,
  });

  const fetchMoreTraces = useCallback(() => {
    const lastPage = data?.pages[data.pages.length - 1];
    const isFirstPage = data?.pages.length === 0;
    const PAGE_SIZE = isFirstPage ? FIRST_PAGE_SIZE : FETCH_MORE_PAGE_SIZE;

    if (!lastPage) {
      fetchNextPage();
      return;
    }

    const isHasMore = lastPage.tracesList.length >= PAGE_SIZE;
    if (isHasMore) {
      fetchNextPage({
        pageParam: {
          cursor: lastPage.cursor,
          minUpdatedAt: lastPage.minUpdatedAt,
        },
      });
      return;
    }

    // Update local state to ensure the next fetch will use the correct minUpdatedAt
    const newMinUpdatedAt = lastPage.minUpdatedAt - TIME_RANGE_LIMIT_IN_SECONDS;
    setLocalMinUpdatedAt(newMinUpdatedAt);
    fetchNextPage({
      pageParam: {
        cursor: lastPage.cursor,
        minUpdatedAt: newMinUpdatedAt,
      },
    });
  }, [data?.pages, fetchNextPage]);

  return {
    isFetching: otherParams.isFetching,
    refetch: otherParams.refetch,
    isSuccess: otherParams.isSuccess,
    fetchNextPage: fetchMoreTraces,
    data: data,
  };
};
