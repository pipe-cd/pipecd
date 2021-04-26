import {
  createAsyncThunk,
  createEntityAdapter,
  createSlice,
  EntityId,
} from "@reduxjs/toolkit";
import {
  Deployment,
  PipelineStage,
  DeploymentStatus,
  StageStatus,
} from "pipe/pkg/app/web/model/deployment_pb";
import * as deploymentsApi from "../api/deployments";
import { fetchCommand, Command, CommandStatus } from "./commands";
import { AppState } from ".";
import { LoadingStatus } from "../types/module";
import { ListDeploymentsRequest } from "pipe/pkg/app/web/api_client/service_pb";
import { ApplicationKind } from "./applications";

export type Stage = Required<PipelineStage.AsObject>;
export type DeploymentStatusKey = keyof typeof DeploymentStatus;

const ITEMS_PER_PAGE = 50;
const FETCH_MORE_ITEMS_PER_PAGE = 30;

export interface DeploymentFilterOptions {
  status?: string;
  kind?: string;
  applicationId?: string;
  envId?: string;
}

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

export const deploymentsAdapter = createEntityAdapter<Deployment.AsObject>({
  sortComparer: (a, b) => b.updatedAt - a.updatedAt,
});

const initialState = deploymentsAdapter.getInitialState<{
  status: LoadingStatus;
  loading: Record<string, boolean>;
  canceling: Record<string, boolean>;
  hasMore: boolean;
}>({
  status: "idle",
  loading: {},
  canceling: {},
  hasMore: true,
});

const selectLastItem = (
  state: typeof initialState
): Deployment.AsObject | undefined => {
  if (state.ids.length === 0) {
    return undefined;
  }
  const lastId = state.ids[state.ids.length - 1];

  return state.entities[lastId];
};

export const fetchDeploymentById = createAsyncThunk<
  Deployment.AsObject,
  string
>("deployments/fetchById", async (deploymentId) => {
  const { deployment } = await deploymentsApi.getDeployment({ deploymentId });
  return deployment as Deployment.AsObject;
});

const convertFilterOptions = (
  options: DeploymentFilterOptions & { maxUpdatedAt: number }
): ListDeploymentsRequest.Options.AsObject => {
  return {
    applicationIdsList: options.applicationId ? [options.applicationId] : [],
    envIdsList: options.envId ? [options.envId] : [],
    kindsList: options.kind
      ? [parseInt(options.kind, 10) as ApplicationKind]
      : [],
    statusesList: options.status
      ? [parseInt(options.status, 10) as DeploymentStatus]
      : [],
    maxUpdatedAt: options.maxUpdatedAt,
  };
};

/**
 * This action will clear old items and add items.
 */
export const fetchDeployments = createAsyncThunk<
  Deployment.AsObject[],
  DeploymentFilterOptions,
  { state: AppState }
>("deployments/fetchList", async (options) => {
  const { deploymentsList } = await deploymentsApi.getDeployments({
    options: convertFilterOptions({ ...options, maxUpdatedAt: 0 }),
    pageSize: ITEMS_PER_PAGE,
    cursor: "",
  });
  return (deploymentsList as Deployment.AsObject[]) || [];
});

/**
 * This action will add items to current state.
 */
export const fetchMoreDeployments = createAsyncThunk<
  Deployment.AsObject[],
  DeploymentFilterOptions,
  { state: AppState }
>("deployments/fetchMoreList", async (options, thunkAPI) => {
  const { deployments } = thunkAPI.getState();
  const lastItem = selectLastItem(deployments);
  const maxUpdatedAt = lastItem ? lastItem.updatedAt : 0;

  const { deploymentsList } = await deploymentsApi.getDeployments({
    options: convertFilterOptions({ ...options, maxUpdatedAt }),
    pageSize: FETCH_MORE_ITEMS_PER_PAGE,
    // TODO: Support cursor passed from frontend.
    cursor: "",
  });
  return (deploymentsList as Deployment.AsObject[]) || [];
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
        state.status = "loading";
        state.hasMore = true;
      })
      .addCase(fetchDeployments.fulfilled, (state, action) => {
        state.status = "succeeded";
        deploymentsAdapter.removeAll(state);
        if (action.payload.length > 0) {
          deploymentsAdapter.upsertMany(state, action.payload);
        }
        if (action.payload.length < ITEMS_PER_PAGE) {
          state.hasMore = false;
        }
      })
      .addCase(fetchDeployments.rejected, (state) => {
        state.status = "failed";
      })
      .addCase(fetchMoreDeployments.pending, (state) => {
        state.status = "loading";
      })
      .addCase(fetchMoreDeployments.fulfilled, (state, action) => {
        state.status = "succeeded";
        deploymentsAdapter.upsertMany(state, action.payload);
        if (action.payload.length < FETCH_MORE_ITEMS_PER_PAGE) {
          state.hasMore = false;
        }
      })
      .addCase(fetchMoreDeployments.rejected, (state) => {
        state.status = "failed";
      })

      .addCase(cancelDeployment.pending, (state, action) => {
        state.canceling[action.meta.arg.deploymentId] = true;
      })
      .addCase(fetchCommand.fulfilled, (state, action) => {
        if (
          action.payload.type === Command.Type.CANCEL_DEPLOYMENT &&
          action.payload.status !== CommandStatus.COMMAND_NOT_HANDLED_YET
        ) {
          state.canceling[action.payload.deploymentId] = false;
        }
      });
  },
});

export const {
  selectById,
  selectAll,
  selectEntities,
  selectIds,
} = deploymentsAdapter.getSelectors();

export const selectDeploymentIsCanceling = (id: EntityId) => (
  state: AppState
): boolean => (id ? state.deployments.canceling[id] : false);

export {
  Deployment,
  DeploymentStatus,
  StageStatus,
  SyncStrategy,
  PipelineStage,
} from "pipe/pkg/app/web/model/deployment_pb";
