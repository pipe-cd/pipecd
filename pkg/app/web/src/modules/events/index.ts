import {
  createAsyncThunk,
  createEntityAdapter,
  createSlice,
} from "@reduxjs/toolkit";
import { Event, EventStatus } from "pipe/pkg/app/web/model/event_pb";
import type { AppState } from "~/store";
import { LoadingStatus } from "~/types/module";

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

/**
 * This action will clear old items and add items.
 */
export const fetchEvents = createAsyncThunk<
  { events: Event.AsObject[]; cursor: string },
  EventFilterOptions,
  { state: AppState }
>(`${MODULE_NAME}/fetchList`, async () => {
  // TODO: Call ListEvents to fetch events
  return {
    events: [
      {
        id: "964a6694-bf3e-4c82-addf-7c8b140cf958",
        name: "push-image",
        data: "v0.2.1",
        projectId: "project-id",
        labelsMap: new Array<[string, string]>(["app", "foo"], ["env", "dev"]),
        eventKey: "push-image:app:foo:env:dev",
        handled: false,
        status: EventStatus.EVENT_NOT_HANDLED,
        statusDescription: "It is going to be replaced by v0.2.1",
        handledAt: 0,
        createdAt: 1642574083,
        updatedAt: 1642574083,
      },
      {
        id: "3d681c85-ab28-458b-9c14-0ddf77634b75",
        name: "helm-release",
        data: "v0.1.0",
        projectId: "project-id",
        labelsMap: new Array<[string, string]>(["app", "bar"]),
        eventKey: "helm-release:app:bar",
        handled: true,
        status: EventStatus.EVENT_SUCCESS,
        statusDescription:
          "Successfully updated 2 files in the repo-1 repository",
        handledAt: 1642574082,
        createdAt: 1642574073,
        updatedAt: 1642574073,
      },
      {
        id: "2a681c85-ab28-458b-9c14-0ddf77634b75",
        name: "helm-release",
        data: "v0.1.0",
        projectId: "project-id",
        labelsMap: new Array<[string, string]>(["app", "bar"]),
        eventKey: "helm-release:app:bar",
        handled: true,
        status: EventStatus.EVENT_FAILURE,
        statusDescription:
          "Failed to change files: failed to get value at $.spec.template.spec.containers[0]version in /path/to/repo/foo/deployment.yaml",
        handledAt: 1642143308,
        createdAt: 1642143305,
        updatedAt: 1642143305,
      },
      {
        id: "8e681c85-ab28-458b-9c14-0ddf77634b75",
        name: "helm-release",
        data: "v0.1.0",
        projectId: "project-id",
        labelsMap: new Array<[string, string]>(["app", "bar"]),
        eventKey: "helm-release:app:bar",
        handled: false,
        status: EventStatus.EVENT_OUTDATED,
        statusDescription: "The new event has been created",
        handledAt: 1642142208,
        createdAt: 1642142205,
        updatedAt: 1642142205,
      },
    ] as Event.AsObject[],
    cursor: "",
  };
});

/**
 * This action will add items to current state.
 */
export const fetchMoreEvents = createAsyncThunk<
  { events: Event.AsObject[]; cursor: string },
  EventFilterOptions,
  { state: AppState }
>(`${MODULE_NAME}/fetchMoreList`, async () => {
  // TODO: Call ListEvents to fetch more events
  return {
    events: [
      {
        id: "a84b1c85-ab28-458b-9c14-0ddf77634b75",
        name: "helm-release",
        data: "v0.1.0",
        projectId: "project-id",
        labelsMap: new Array<[string, string]>(["app", "bar"]),
        eventKey: "helm-release:app:bar",
        handled: true,
        status: EventStatus.EVENT_FAILURE,
        statusDescription:
          "Failed to change files: failed to get value at $.spec.template.spec.containers[0]version in /path/to/repo/foo/deployment.yaml",
        handledAt: 1642111308,
        createdAt: 1642111305,
        updatedAt: 1642111305,
      },
    ] as Event.AsObject[],
    cursor: "",
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
          state.minUpdatedAt =
            events[events.length - 1].updatedAt - TIME_RANGE_LIMIT_IN_SECONDS;
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

export { Event, EventStatus } from "pipe/pkg/app/web/model/event_pb";
