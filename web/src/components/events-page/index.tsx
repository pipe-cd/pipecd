import {
  Box,
  Button,
  CircularProgress,
  Divider,
  List,
  makeStyles,
  Toolbar,
  Typography,
} from "@material-ui/core";
import CloseIcon from "@material-ui/icons/Close";
import FilterIcon from "@material-ui/icons/FilterList";
import RefreshIcon from "@material-ui/icons/Refresh";
import dayjs from "dayjs";
import { FC, useCallback, useEffect, useRef, useState } from "react";
import { useInView } from "react-intersection-observer";
import { useNavigate } from "react-router-dom";
import { PAGE_PATH_EVENTS } from "~/constants/path";
import {
  UI_TEXT_FILTER,
  UI_TEXT_HIDE_FILTER,
  UI_TEXT_REFRESH,
  UI_TEXT_MORE,
} from "~/constants/ui-text";
import {
  useAppDispatch,
  useAppSelector,
  useShallowEqualSelector,
} from "~/hooks/redux";
import { fetchApplications } from "~/modules/applications";
import {
  Event,
  EventFilterOptions,
  fetchEvents,
  fetchMoreEvents,
  selectById as selectEventById,
  selectIds as selectEventIds,
} from "~/modules/events";
import { useStyles as useButtonStyles } from "~/styles/button";
import {
  stringifySearchParams,
  useSearchParams,
  arrayFormat,
} from "~/utils/search-params";
import { EventFilter } from "./event-filter";
import { EventItem } from "./event-item";

const useStyles = makeStyles((theme) => ({
  eventLists: {
    listStyle: "none",
    padding: theme.spacing(3),
    paddingTop: 0,
    margin: 0,
    flex: 1,
    overflowY: "scroll",
  },
  date: {
    marginTop: theme.spacing(2),
    marginBottom: theme.spacing(2),
  },
}));

const sortComp = (a: string | number, b: string | number): number => {
  return dayjs(b).valueOf() - dayjs(a).valueOf();
};

function filterUndefined<TValue>(value: TValue | undefined): value is TValue {
  return value !== undefined;
}

const useGroupedEvents = (): Record<string, Event.AsObject[]> => {
  const events = useShallowEqualSelector<Event.AsObject[]>((state) =>
    selectEventIds(state.events)
      .map((id) => selectEventById(state.events, id))
      .filter(filterUndefined)
  );

  const result: Record<string, Event.AsObject[]> = {};

  events.forEach((event) => {
    const dateStr = dayjs(event.createdAt * 1000).format("YYYY/MM/DD");
    if (!result[dateStr]) {
      result[dateStr] = [];
    }
    result[dateStr].push(event);
  });

  return result;
};

export const EventIndexPage: FC = () => {
  const classes = useStyles();
  const buttonClasses = useButtonStyles();
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const listRef = useRef(null);
  const status = useAppSelector((state) => state.events.status);
  const hasMore = useAppSelector((state) => state.events.hasMore);
  const groupedEvents = useGroupedEvents();
  const filterOptions = useSearchParams();
  const [openFilter, setOpenFilter] = useState(true);
  const [ref, inView] = useInView({
    rootMargin: "400px",
    root: listRef.current,
  });

  const isLoading = status === "loading";

  useEffect(() => {
    dispatch(fetchApplications());
  }, [dispatch]);

  useEffect(() => {
    dispatch(fetchEvents(filterOptions));
  }, [dispatch, filterOptions]);

  useEffect(() => {
    if (inView && hasMore && isLoading === false) {
      dispatch(fetchMoreEvents(filterOptions));
    }
  }, [dispatch, inView, hasMore, isLoading, filterOptions]);

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
    dispatch(fetchEvents(filterOptions));
  }, [dispatch, filterOptions]);

  const handleMoreClick = useCallback(() => {
    dispatch(fetchMoreEvents(filterOptions));
  }, [dispatch, filterOptions]);

  const dates = Object.keys(groupedEvents).sort(sortComp);

  return (
    <Box display="flex" overflow="hidden" flex={1} flexDirection="column">
      <Toolbar variant="dense">
        <Box flexGrow={1} />
        <Button
          color="primary"
          startIcon={<RefreshIcon />}
          onClick={handleRefreshClick}
          disabled={isLoading}
        >
          {UI_TEXT_REFRESH}
          {isLoading && (
            <CircularProgress size={24} className={buttonClasses.progress} />
          )}
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
      <Box display="flex" overflow="hidden" flex={1}>
        <ol className={classes.eventLists} ref={listRef}>
          {dates.length === 0 &&
            (isLoading ? (
              <Box display="flex" justifyContent="center" mt={3}>
                <CircularProgress />
              </Box>
            ) : (
              <Box display="flex" justifyContent="center" mt={3}>
                <Typography>No events</Typography>
              </Box>
            ))}
          {dates.map((date) => (
            <li key={date}>
              <Typography variant="subtitle1" className={classes.date}>
                {date}
              </Typography>
              <List>
                {groupedEvents[date]
                  .sort((a, b) => sortComp(a.createdAt, b.createdAt))
                  .map((event) => (
                    <EventItem id={event.id} key={`event-item-${event.id}`} />
                  ))}
              </List>
            </li>
          ))}
          {status === "succeeded" && <div ref={ref} />}
          {!hasMore && (
            <Button
              color="primary"
              variant="outlined"
              size="large"
              fullWidth
              onClick={handleMoreClick}
              disabled={isLoading}
            >
              {UI_TEXT_MORE}
              {isLoading && (
                <CircularProgress
                  size={24}
                  className={buttonClasses.progress}
                />
              )}
            </Button>
          )}
        </ol>
        {openFilter && (
          <EventFilter
            options={filterOptions}
            onChange={handleFilterChange}
            onClear={handleFilterClear}
          />
        )}
      </Box>
    </Box>
  );
};
