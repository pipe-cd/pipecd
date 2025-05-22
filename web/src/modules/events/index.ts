import {
  createAsyncThunk,
  createEntityAdapter,
  createSlice,
} from "@reduxjs/toolkit";
import { Event, EventStatus } from "pipecd/web/model/event_pb";
import * as eventsApi from "~/api/events";
import type { AppState } from "~/store";
import { LoadingStatus } from "~/types/module";
import { ListEventsRequest } from "pipecd/web/api_client/service_pb";

export type EventStatusKey = keyof typeof EventStatus;

// 30 days
const TIME_RANGE_LIMIT_IN_SECONDS = 2592000;
const ITEMS_PER_PAGE = 50;
const FETCH_MORE_ITEMS_PER_PAGE = 30;

const MODULE_NAME = "events";

export interface EventFilterOptions {
  name?: string;
  status?: string;
  // Suppose to be like ["key-1:value-1"]
  // sindresorhus/query-string doesn't support multidimensional arrays, that's why the format is a bit tricky.
  labels?: Array<string>;
}

export const eventsAdapter = createEntityAdapter<Event.AsObject>({
  sortComparer: (a, b) => b.updatedAt - a.updatedAt,
});

const initialState = eventsAdapter.getInitialState<{
  status: LoadingStatus;
  loading: Record<string, boolean>;
  hasMore: boolean;
  cursor: string;
  minUpdatedAt: number;
}>({
  status: "idle",
  loading: {},
  hasMore: true,
  cursor: "",
  minUpdatedAt: Math.round(Date.now() / 1000 - TIME_RANGE_LIMIT_IN_SECONDS),
});

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

/**
 * This action will clear old items and add items.
 */
export const fetchEvents = createAsyncThunk<
  { events: Event.AsObject[]; cursor: string },
  EventFilterOptions,
  { state: AppState }
>(`${MODULE_NAME}/fetchList`, async (options, thunkAPI) => {
  const { events } = thunkAPI.getState();
  const { eventsList, cursor } = await eventsApi.getEvents({
    options: convertFilterOptions({ ...options }),
    pageSize: ITEMS_PER_PAGE,
    cursor: "",
    pageMinUpdatedAt: events.minUpdatedAt,
  });

  return {
    events: (eventsList as Event.AsObject[]) || [],
    cursor,
  };
});

/**
 * This action will add items to current state.
 */
export const fetchMoreEvents = createAsyncThunk<
  { events: Event.AsObject[]; cursor: string },
  EventFilterOptions,
  { state: AppState }
>(`${MODULE_NAME}/fetchMoreList`, async (options, thunkAPI) => {
  const { events } = thunkAPI.getState();
  const { eventsList, cursor } = await eventsApi.getEvents({
    options: convertFilterOptions({ ...options }),
    pageSize: FETCH_MORE_ITEMS_PER_PAGE,
    cursor: events.cursor,
    pageMinUpdatedAt: events.minUpdatedAt,
  });

  return {
    events: (eventsList as Event.AsObject[]) || [],
    cursor,
  };
});

export const eventsSlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchEvents.pending, (state) => {
        state.status = "loading";
        state.hasMore = true;
        state.cursor = "";
      })
      .addCase(fetchEvents.fulfilled, (state, action) => {
        state.status = "succeeded";
        eventsAdapter.removeAll(state);
        if (action.payload.events.length > 0) {
          eventsAdapter.upsertMany(state, action.payload.events);
        }
        if (action.payload.events.length < ITEMS_PER_PAGE) {
          state.hasMore = false;
        }
        state.cursor = action.payload.cursor;
      })
      .addCase(fetchEvents.rejected, (state) => {
        state.status = "failed";
      })
      .addCase(fetchMoreEvents.pending, (state) => {
        state.status = "loading";
      })
      .addCase(fetchMoreEvents.fulfilled, (state, action) => {
        state.status = "succeeded";
        eventsAdapter.upsertMany(state, action.payload.events);
        const events = action.payload.events;
        if (events.length < FETCH_MORE_ITEMS_PER_PAGE) {
          state.hasMore = false;
          state.minUpdatedAt = state.minUpdatedAt - TIME_RANGE_LIMIT_IN_SECONDS;
        } else {
          state.hasMore = true;
        }
        state.cursor = action.payload.cursor;
      })
      .addCase(fetchMoreEvents.rejected, (state) => {
        state.status = "failed";
      });
  },
});

export const {
  selectById,
  selectAll,
  selectEntities,
  selectIds,
} = eventsAdapter.getSelectors();

export { Event, EventStatus } from "pipecd/web/model/event_pb";
