import {
  useInfiniteQuery,
  UseInfiniteQueryResult,
} from "@tanstack/react-query";
import { ListDeploymentsRequest } from "pipecd/web/api_client/service_pb";
import { ApplicationKind } from "~~/model/common_pb";
import { Deployment, DeploymentStatus } from "pipecd/web/model/deployment_pb";
import * as deploymentsApi from "~/api/deployments";

// 30 days
const TIME_RANGE_LIMIT_IN_SECONDS = 2592000;
const ITEMS_PER_PAGE = 50;
const FETCH_MORE_ITEMS_PER_PAGE = 30;

// TODO DONE
export interface DeploymentFilterOptions {
  status?: string;
  kind?: string;
  applicationId?: string;
  applicationName?: string;
  // Suppose to be like ["key-1:value-1"]
  // sindresorhus/query-string doesn't support multidimensional arrays, that's why the format is a bit tricky.
  labels?: Array<string>;
}

// TODO DONE
const convertFilterOptions = (
  options: DeploymentFilterOptions
): ListDeploymentsRequest.Options.AsObject => {
  const labels = new Array<[string, string]>();
  if (options.labels) {
    for (const label of options.labels) {
      const pair = label.split(":");
      if (pair.length === 2) labels.push([pair[0], pair[1]]);
    }
  }
  return {
    applicationName: options.applicationName ?? "",
    applicationIdsList: options.applicationId ? [options.applicationId] : [],
    kindsList: options.kind
      ? [parseInt(options.kind, 10) as ApplicationKind]
      : [],
    statusesList: options.status
      ? [parseInt(options.status, 10) as DeploymentStatus]
      : [],
    labelsMap: labels,
  };
};

export const useGetDeploymentsInfinite = (
  options: DeploymentFilterOptions
): UseInfiniteQueryResult<{
  deploymentsList: Deployment.AsObject[];
  cursor: string;
  minUpdatedAt: number;
  hasMore: boolean;
}> => {
  // TODO
  return useInfiniteQuery({
    queryKey: ["deployment", options],
    queryFn: async ({ pageParam }) => {
      const isFirstPage = !pageParam;
      const minUpdatedAt =
        pageParam?.minUpdatedAt ??
        Math.round(Date.now() / 1000 - TIME_RANGE_LIMIT_IN_SECONDS);

      const pageSize = isFirstPage ? ITEMS_PER_PAGE : FETCH_MORE_ITEMS_PER_PAGE;

      const { deploymentsList, cursor } = await deploymentsApi.getDeployments({
        options: convertFilterOptions({ ...options }),
        pageSize: pageSize,
        cursor: pageParam?.cursor || "",
        pageMinUpdatedAt: minUpdatedAt,
      });

      return {
        deploymentsList,
        cursor: cursor || pageParam?.cursor || "",
        minUpdatedAt,
        hasMore: deploymentsList.length >= pageSize,
      };
    },
    getNextPageParam: (lastPage, allPages) => {
      const isFirstPage = allPages.length === 0;

      const ITEMS_PER_PAGE_TO_USE = isFirstPage
        ? ITEMS_PER_PAGE
        : FETCH_MORE_ITEMS_PER_PAGE;

      const isHasMore =
        lastPage.deploymentsList.length >= ITEMS_PER_PAGE_TO_USE;

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
