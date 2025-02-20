import {
  createAsyncThunk,
  createEntityAdapter,
  createSlice,
} from "@reduxjs/toolkit";
import { LoadingStatus } from "~/types/module";
import * as deploymentTracesApi from "~/api/deploymentTraces";
import { AppState } from "~/store";
import {
  ListDeploymentTracesRequest,
  ListDeploymentTracesResponse,
} from "~~/api_client/service_pb";

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

export const deploymentTraceAdapter = createEntityAdapter<
  ListDeploymentTracesResponse.DeploymentTraceRes.AsObject
>({
  selectId: (trace) => trace.trace?.id as string,
  sortComparer: (a, b) => {
    if (!b.trace?.updatedAt) return 0;
    if (!a.trace?.updatedAt) return 0;
    return b.trace?.updatedAt - a.trace?.updatedAt;
  },
});

const initialState = deploymentTraceAdapter.getInitialState<{
  status: LoadingStatus;
  loading: Record<string, boolean>;
  hasMore: boolean;
  cursor: string;
  minUpdatedAt: number;
  skippable: Record<string, boolean | undefined>;
}>({
  status: "idle",
  loading: {},
  hasMore: true,
  cursor: "",
  minUpdatedAt: Math.round(Date.now() / 1000 - TIME_RANGE_LIMIT_IN_SECONDS),
  skippable: {},
});

export const fetchDeploymentTraces = createAsyncThunk<
  {
    tracesList?: ListDeploymentTracesResponse.DeploymentTraceRes.AsObject[];
    cursor: string;
  },
  DeploymentTraceFilterOptions,
  { state: AppState }
>("deploymentTrace/fetchList", async (options, thunkAPI) => {
  const { deploymentTrace } = thunkAPI.getState();

  const response = await deploymentTracesApi.getDeploymentTraces({
    options: convertFilterOptions(options),
    pageSize: ITEMS_PER_PAGE,
    cursor: "",
    pageMinUpdatedAt: deploymentTrace.minUpdatedAt,
  });

  return response;
});

export const fetchMoreDeploymentTraces = createAsyncThunk<
  {
    tracesList: ListDeploymentTracesResponse.DeploymentTraceRes.AsObject[];
    cursor: string;
  },
  DeploymentTraceFilterOptions,
  { state: AppState }
>("deploymentTrace/fetchMoreList", async (options, thunkAPI) => {
  const { deployments } = thunkAPI.getState();

  const response = await deploymentTracesApi.getDeploymentTraces({
    options: convertFilterOptions(options),
    pageSize: FETCH_MORE_ITEMS_PER_PAGE,
    cursor: deployments.cursor,
    pageMinUpdatedAt: deployments.minUpdatedAt,
  });

  return response;
});

export const deploymentTraceSlice = createSlice({
  name: "deploymentTrace",
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchDeploymentTraces.pending, (state) => {
        state.status = "loading";
        state.hasMore = true;
        state.cursor = "";
      })
      .addCase(fetchDeploymentTraces.fulfilled, (state, action) => {
        state.status = "succeeded";
        deploymentTraceAdapter.removeAll(state);
        if (action.payload.tracesList) {
          if (action.payload.tracesList?.length > 0) {
            deploymentTraceAdapter.upsertMany(state, action.payload.tracesList);
          }
          if (action.payload.tracesList?.length < ITEMS_PER_PAGE) {
            state.hasMore = false;
          }
        }
        state.cursor = action.payload.cursor;
      })
      .addCase(fetchDeploymentTraces.rejected, (state) => {
        state.status = "failed";
      })
      .addCase(fetchMoreDeploymentTraces.pending, (state) => {
        state.status = "loading";
      })
      .addCase(fetchMoreDeploymentTraces.fulfilled, (state, action) => {
        state.status = "succeeded";
        if (action.payload.tracesList.length > 0) {
          deploymentTraceAdapter.upsertMany(state, action.payload.tracesList);
        }
        if (action.payload.tracesList.length < FETCH_MORE_ITEMS_PER_PAGE) {
          state.hasMore = false;
          state.minUpdatedAt = state.minUpdatedAt - TIME_RANGE_LIMIT_IN_SECONDS;
        } else {
          state.hasMore = true;
        }
        state.cursor = action.payload.cursor;
      })
      .addCase(fetchMoreDeploymentTraces.rejected, (state) => {
        state.status = "failed";
      });
  },
});

export const {
  selectById,
  selectAll,
  selectEntities,
  selectIds,
} = deploymentTraceAdapter.getSelectors();
