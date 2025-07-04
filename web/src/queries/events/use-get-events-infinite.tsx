import {
  useInfiniteQuery,
  UseInfiniteQueryResult,
} from "@tanstack/react-query";
import { ListEventsRequest } from "~~/api_client/service_pb";
import { Event, EventStatus } from "pipecd/web/model/event_pb";
import * as eventsApi from "~/api/events";

// 30 days
const TIME_RANGE_LIMIT_IN_SECONDS = 2592000;
const ITEMS_PER_PAGE = 50;
const FETCH_MORE_ITEMS_PER_PAGE = 30;

export type EventStatusKey = keyof typeof EventStatus;

export type EventFilterOptions = {
  name?: string;
  status?: string;
  // Suppose to be like ["key-1:value-1"]
  // sindresorhus/query-string doesn't support multidimensional arrays, that's why the format is a bit tricky.
  labels?: Array<string>;
};

const convertFilterOptions = (
  options: EventFilterOptions
): ListEventsRequest.Options.AsObject => {
  const labels = new Array<[string, string]>();
  if (options.labels) {
    for (const label of options.labels) {
      const pair = label.split(":");
      if (pair.length === 2) labels.push([pair[0], pair[1]]);
    }
  }
  return {
    name: options.name ?? "",
    statusesList: options.status
      ? [parseInt(options.status, 10) as EventStatus]
      : [],
    labelsMap: labels,
  };
};

export const useGetEventsInfinite = (
  options: EventFilterOptions
): UseInfiniteQueryResult<{
  eventsList: Event.AsObject[];
  cursor: string;
  minUpdatedAt: number;
  hasMore: boolean;
}> => {
  return useInfiniteQuery({
    queryKey: ["events", options],
    queryFn: async ({ pageParam }) => {
      const isFirstPage = !pageParam;
      const minUpdatedAt =
        pageParam?.minUpdatedAt ??
        Math.round(Date.now() / 1000 - TIME_RANGE_LIMIT_IN_SECONDS);

      const pageSize = isFirstPage ? ITEMS_PER_PAGE : FETCH_MORE_ITEMS_PER_PAGE;

      const { eventsList, cursor } = await eventsApi.getEvents({
        options: convertFilterOptions({ ...options }),
        pageSize: pageSize,
        cursor: pageParam?.cursor || "",
        pageMinUpdatedAt: minUpdatedAt,
      });

      return {
        eventsList,
        cursor: cursor || pageParam?.cursor || "",
        minUpdatedAt,
        hasMore: eventsList.length >= pageSize,
      };
    },
    getNextPageParam: (lastPage, allPages) => {
      const isFirstPage = allPages.length === 0;

      const ITEMS_PER_PAGE_TO_USE = isFirstPage
        ? ITEMS_PER_PAGE
        : FETCH_MORE_ITEMS_PER_PAGE;

      const isHasMore = lastPage.eventsList.length >= ITEMS_PER_PAGE_TO_USE;

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
