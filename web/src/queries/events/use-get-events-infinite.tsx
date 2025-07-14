import {
  useInfiniteQuery,
  UseInfiniteQueryResult,
} from "@tanstack/react-query";
import { ListEventsRequest } from "~~/api_client/service_pb";
import { Event, EventStatus } from "pipecd/web/model/event_pb";
import * as eventsApi from "~/api/events";
import { useCallback, useState } from "react";

// 30 days
const TIME_RANGE_LIMIT_IN_SECONDS = 2592000;
const FIRST_PAGE_SIZE = 50;
const FETCH_MORE_PAGE_SIZE = 30;

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

type QueryType = UseInfiniteQueryResult<{
  eventsList: Event.AsObject[];
  cursor: string;
  minUpdatedAt: number;
  hasMore: boolean;
}>;

export const useGetEventsInfinite = (
  options: EventFilterOptions
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
    queryKey: ["events", options],
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
    refetchOnWindowFocus: false,
  });

  const fetchMoreEvents = useCallback(() => {
    const lastPage = data?.pages[data.pages.length - 1];
    const isFirstPage = data?.pages.length === 0;
    const PAGE_SIZE = isFirstPage ? FIRST_PAGE_SIZE : FETCH_MORE_PAGE_SIZE;

    if (!lastPage) {
      fetchNextPage();
      return;
    }

    const isHasMore = lastPage.eventsList.length >= PAGE_SIZE;
    if (isHasMore) {
      fetchNextPage();
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
    fetchNextPage: fetchMoreEvents,
    data: data,
  };
};
