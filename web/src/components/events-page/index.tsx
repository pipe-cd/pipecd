import {
  Box,
  Button,
  CircularProgress,
  Divider,
  List,
  Toolbar,
  Typography,
} from "@mui/material";
import CloseIcon from "@mui/icons-material/Close";
import FilterIcon from "@mui/icons-material/FilterList";
import RefreshIcon from "@mui/icons-material/Refresh";
import dayjs from "dayjs";
import { FC, useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useInView } from "react-intersection-observer";
import { useNavigate } from "react-router-dom";
import { PAGE_PATH_EVENTS } from "~/constants/path";
import {
  UI_TEXT_FILTER,
  UI_TEXT_HIDE_FILTER,
  UI_TEXT_REFRESH,
  UI_TEXT_MORE,
} from "~/constants/ui-text";

import { Event } from "pipecd/web/model/event_pb";
import { SpinnerIcon } from "~/styles/button";
import {
  stringifySearchParams,
  useSearchParams,
  arrayFormat,
} from "~/utils/search-params";
import { EventFilter } from "./event-filter";
import { EventItem } from "./event-item";
import {
  EventFilterOptions,
  useGetEventsInfinite,
} from "~/queries/events/use-get-events-infinite";

const sortComp = (a: string | number, b: string | number): number => {
  return dayjs(b).valueOf() - dayjs(a).valueOf();
};

const useNewGroupedEvents = (
  events: Event.AsObject[]
): Record<string, Event.AsObject[]> => {
  return useMemo(() => {
    const result: Record<string, Event.AsObject[]> = {};

    events.forEach((event) => {
      const dateStr = dayjs(event.createdAt * 1000).format("YYYY/MM/DD");
      if (!result[dateStr]) {
        result[dateStr] = [];
      }
      if (!result[dateStr].some((e) => e.id === event.id)) {
        result[dateStr].push(event);
      }
    });

    return result;
  }, [events]);
};

export const EventIndexPage: FC = () => {
  const navigate = useNavigate();
  const listRef = useRef(null);
  const filterOptions = useSearchParams();
  const [openFilter, setOpenFilter] = useState(true);
  const [ref, inView] = useInView({
    rootMargin: "400px",
    root: listRef.current,
  });

  const {
    data: eventsData,
    isFetching,
    fetchNextPage: fetchMoreEvents,
    refetch: refetchEvents,
    isSuccess,
  } = useGetEventsInfinite(filterOptions);

  const eventsList = useMemo(() => {
    return eventsData?.pages.flatMap((item) => item.eventsList) || [];
  }, [eventsData]);

  const hasMore = useMemo(() => {
    if (!eventsData || eventsData.pages.length === 0) {
      return false;
    }
    const lastIndex = eventsData?.pages.length - 1;
    return eventsData?.pages?.[lastIndex]?.hasMore || false;
  }, [eventsData]);

  const groupedEvents = useNewGroupedEvents(eventsList || []);

  useEffect(() => {
    if (inView && hasMore && isFetching === false) {
      fetchMoreEvents();
    }
  }, [inView, isFetching, filterOptions, fetchMoreEvents, hasMore]);

  // filter handlers
  const handleFilterChange = useCallback(
    (options: EventFilterOptions) => {
      navigate(
        `${PAGE_PATH_EVENTS}?${stringifySearchParams(
          { ...options },
          { arrayFormat: arrayFormat }
        )}`,
        { replace: true }
      );
    },
    [navigate]
  );
  const handleFilterClear = useCallback(() => {
    navigate(PAGE_PATH_EVENTS, { replace: true });
  }, [navigate]);

  const handleRefreshClick = useCallback(() => {
    refetchEvents();
  }, [refetchEvents]);

  const handleMoreClick = useCallback(() => {
    fetchMoreEvents();
  }, [fetchMoreEvents]);

  const dates = Object.keys(groupedEvents).sort(sortComp);

  return (
    <Box
      sx={{
        display: "flex",
        overflow: "hidden",
        flex: 1,
        flexDirection: "column",
      }}
    >
      <Toolbar variant="dense">
        <Box
          sx={{
            flexGrow: 1,
          }}
        />
        <Button
          color="primary"
          startIcon={<RefreshIcon />}
          onClick={handleRefreshClick}
          disabled={isFetching}
        >
          {UI_TEXT_REFRESH}
          {isFetching && <SpinnerIcon />}
        </Button>
        <Button
          color="primary"
          startIcon={openFilter ? <CloseIcon /> : <FilterIcon />}
          onClick={() => setOpenFilter(!openFilter)}
        >
          {openFilter ? UI_TEXT_HIDE_FILTER : UI_TEXT_FILTER}
        </Button>
      </Toolbar>
      <Divider />
      <Box
        sx={{
          display: "flex",
          overflow: "hidden",
          flex: 1,
        }}
      >
        <Box
          component={"ol"}
          sx={{
            listStyle: "none",
            padding: 3,
            paddingTop: 0,
            margin: 0,
            flex: 1,
            overflowY: "scroll",
          }}
          ref={listRef}
        >
          {dates.length === 0 && (
            <Box
              sx={{
                display: "flex",
                justifyContent: "center",
                mt: 3,
              }}
            >
              {isFetching ? (
                <CircularProgress />
              ) : (
                <Typography>No events</Typography>
              )}
            </Box>
          )}
          {dates.map((date) => (
            <li key={date}>
              <Typography
                variant="subtitle1"
                sx={{
                  mt: 2,
                  mb: 2,
                }}
              >
                {date}
              </Typography>
              <List>
                {groupedEvents[date]
                  .sort((a, b) => sortComp(a.createdAt, b.createdAt))
                  .map((event) => {
                    return (
                      <EventItem event={event} key={`event-item-${event.id}`} />
                    );
                  })}
              </List>
            </li>
          ))}
          {isSuccess && <div ref={ref} />}
          {!eventsData?.pages?.[0]?.hasMore && (
            <Button
              color="primary"
              variant="outlined"
              size="large"
              fullWidth
              onClick={handleMoreClick}
              disabled={isFetching}
            >
              {UI_TEXT_MORE}
              {isFetching && <SpinnerIcon />}
            </Button>
          )}
        </Box>
        {openFilter && (
          <EventFilter
            events={eventsList}
            options={filterOptions}
            onChange={handleFilterChange}
            onClear={handleFilterClear}
          />
        )}
      </Box>
    </Box>
  );
};
