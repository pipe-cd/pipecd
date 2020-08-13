import {
  createAsyncThunk,
  createEntityAdapter,
  createSlice,
} from "@reduxjs/toolkit";
import {
  Deployment as DeploymentModel,
  PipelineStage,
  DeploymentStatus,
  StageStatus,
} from "pipe/pkg/app/web/model/deployment_pb";
import * as deploymentsApi from "../api/deployments";
import { fetchCommand, CommandModel, CommandStatus } from "./commands";
import { AppState } from ".";

export type Deployment = Required<DeploymentModel.AsObject>;
export type Stage = Required<PipelineStage.AsObject>;
export type DeploymentStatusKey = keyof typeof DeploymentStatus;

const ITEMS_PER_PAGE = 50;
const FETCH_MORE_ITEMS_PER_PAGE = 30;

export const isDeploymentRunning = (
  status: DeploymentStatus | undefined
): boolean => {
  if (status === undefined) {
    return false;
  }

  switch (status) {
    case DeploymentStatus.DEPLOYMENT_PENDING:
    case DeploymentStatus.DEPLOYMENT_PLANNED:
    case DeploymentStatus.DEPLOYMENT_ROLLING_BACK:
    case DeploymentStatus.DEPLOYMENT_RUNNING:
      return true;
    case DeploymentStatus.DEPLOYMENT_CANCELLED:
    case DeploymentStatus.DEPLOYMENT_FAILURE:
    case DeploymentStatus.DEPLOYMENT_SUCCESS:
      return false;
  }
};

export const isStageRunning = (status: StageStatus): boolean => {
  switch (status) {
    case StageStatus.STAGE_NOT_STARTED_YET:
    case StageStatus.STAGE_RUNNING:
      return true;
    case StageStatus.STAGE_SUCCESS:
    case StageStatus.STAGE_FAILURE:
    case StageStatus.STAGE_CANCELLED:
      return false;
  }
};

export const deploymentsAdapter = createEntityAdapter<Deployment>({
  sortComparer: (a, b) => b.updatedAt - a.updatedAt,
});

export const {
  selectById,
  selectAll,
  selectEntities,
  selectIds,
} = deploymentsAdapter.getSelectors();

const initialState = deploymentsAdapter.getInitialState<{
  loading: Record<string, boolean>;
  canceling: Record<string, boolean>;
  isLoadingItems: boolean;
  isLoadingMoreItems: boolean;
  hasMore: boolean;
}>({
  loading: {},
  canceling: {},
  isLoadingItems: false,
  isLoadingMoreItems: false,
  hasMore: true,
});

const selectLastItem = (state: typeof initialState): Deployment | undefined => {
  if (state.ids.length === 0) {
    return undefined;
  }
  const lastId = state.ids[state.ids.length - 1];

  return state.entities[lastId];
};

export const fetchDeploymentById = createAsyncThunk<Deployment, string>(
  "deployments/fetchById",
  async (deploymentId) => {
    const { deployment } = await deploymentsApi.getDeployment({ deploymentId });
    return deployment as Deployment;
  }
);

/**
 * This action will clear old items and add items.
 */
export const fetchDeployments = createAsyncThunk<
  Deployment[],
  void,
  { state: AppState }
>("deployments/fetchList", async (_, thunkAPI) => {
  const { deploymentFilterOptions } = thunkAPI.getState();
  const { deploymentsList } = await deploymentsApi.getDeployments({
    options: {
      applicationIdsList: deploymentFilterOptions.applicationIds,
      envIdsList: deploymentFilterOptions.envIds,
      kindsList: deploymentFilterOptions.kinds,
      statusesList: deploymentFilterOptions.statuses,
      maxUpdatedAt: 0,
    },
    pageSize: ITEMS_PER_PAGE,
  });
  return (deploymentsList as Deployment[]) || [];
});

/**
 * This action will add items to current state.
 */
export const fetchMoreDeployments = createAsyncThunk<
  Deployment[],
  void,
  { state: AppState }
>("deployments/fetchMoreList", async (_, thunkAPI) => {
  const { deployments, deploymentFilterOptions } = thunkAPI.getState();
  const lastItem = selectLastItem(deployments);
  const maxUpdatedAt = lastItem ? lastItem.updatedAt : 0;
  const { deploymentsList } = await deploymentsApi.getDeployments({
    options: {
      applicationIdsList: deploymentFilterOptions.applicationIds,
      envIdsList: deploymentFilterOptions.envIds,
      kindsList: deploymentFilterOptions.kinds,
      statusesList: deploymentFilterOptions.statuses,
      maxUpdatedAt,
    },
    pageSize: FETCH_MORE_ITEMS_PER_PAGE,
  });
  return (deploymentsList as Deployment[]) || [];
});

export const approveStage = createAsyncThunk<
  void,
  { deploymentId: string; stageId: string }
>("deployments/approve", async (props, thunkAPI) => {
  const { commandId } = await deploymentsApi.approveStage(props);
  await thunkAPI.dispatch(fetchCommand(commandId));
});

export const cancelDeployment = createAsyncThunk<
  void,
  {
    deploymentId: string;
    forceRollback: boolean;
    forceNoRollback: boolean;
  }
>(
  "deployments/cancel",
  async ({ deploymentId, forceRollback, forceNoRollback }, thunkAPI) => {
    const { commandId } = await deploymentsApi.cancelDeployment({
      deploymentId,
      forceRollback,
      forceNoRollback,
    });

    await thunkAPI.dispatch(fetchCommand(commandId));
  }
);

export const deploymentsSlice = createSlice({
  name: "deployments",
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchDeploymentById.pending, (state, action) => {
        state.loading[action.meta.arg] = true;
      })
      .addCase(fetchDeploymentById.fulfilled, (state, action) => {
        state.loading[action.meta.arg] = false;
        if (action.payload) {
          deploymentsAdapter.upsertOne(state, action.payload);
        }
      })
      .addCase(fetchDeploymentById.rejected, (state, action) => {
        state.loading[action.meta.arg] = false;
      })
      .addCase(fetchDeployments.pending, (state) => {
        state.isLoadingItems = true;
        state.hasMore = true;
      })
      .addCase(fetchDeployments.fulfilled, (state, action) => {
        deploymentsAdapter.removeAll(state);
        if (action.payload.length > 0) {
          deploymentsAdapter.upsertMany(state, action.payload);
        }
        state.isLoadingItems = false;
        if (action.payload.length < ITEMS_PER_PAGE) {
          state.hasMore = false;
        }
      })
      .addCase(fetchDeployments.rejected, (state) => {
        state.isLoadingItems = false;
      })
      .addCase(fetchMoreDeployments.pending, (state) => {
        state.isLoadingMoreItems = true;
      })
      .addCase(fetchMoreDeployments.fulfilled, (state, action) => {
        deploymentsAdapter.upsertMany(state, action.payload);
        state.isLoadingMoreItems = false;
        if (action.payload.length < FETCH_MORE_ITEMS_PER_PAGE) {
          state.hasMore = false;
        }
      })
      .addCase(fetchMoreDeployments.rejected, (state) => {
        state.isLoadingMoreItems = false;
      })

      .addCase(cancelDeployment.pending, (state, action) => {
        state.canceling[action.meta.arg.deploymentId] = true;
      })
      .addCase(fetchCommand.fulfilled, (state, action) => {
        if (
          action.payload.type === CommandModel.Type.CANCEL_DEPLOYMENT &&
          action.payload.status !== CommandStatus.COMMAND_NOT_HANDLED_YET
        ) {
          state.canceling[action.payload.deploymentId] = false;
        }
      });
  },
});

export { DeploymentStatus } from "pipe/pkg/app/web/model/deployment_pb";
