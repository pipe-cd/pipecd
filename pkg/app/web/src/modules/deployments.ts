import {
  createAsyncThunk,
  createEntityAdapter,
  createSlice,
} from "@reduxjs/toolkit";
import {
  Deployment as DeploymentModel,
  PipelineStage,
  DeploymentStatus,
} from "pipe/pkg/app/web/model/deployment_pb";
import * as deploymentsApi from "../api/deployments";
import { fetchCommand, CommandModel, CommandStatus } from "./commands";
import { ApplicationKind } from "./applications";

export type Deployment = Required<DeploymentModel.AsObject>;
export type Stage = Required<PipelineStage.AsObject>;
export type DeploymentStatusKey = keyof typeof DeploymentStatus;

export const isDeploymentRunning = (status: DeploymentStatus): boolean => {
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

export const deploymentsAdapter = createEntityAdapter<Deployment>({});

export const { selectById, selectAll } = deploymentsAdapter.getSelectors();

export const fetchDeploymentById = createAsyncThunk<Deployment, string>(
  "deployments/fetchById",
  async (deploymentId) => {
    const { deployment } = await deploymentsApi.getDeployment({ deploymentId });
    return deployment as Deployment;
  }
);

export const fetchDeployments = createAsyncThunk<
  Deployment[],
  {
    statusesList: DeploymentStatus[];
    kindsList: ApplicationKind[];
    applicationIdsList: string[];
    envIdsList: string[];
  }
>("deployments/fetchList", async (options) => {
  const { deploymentsList } = await deploymentsApi.getDeployments({ options });
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
    withoutRollback: boolean;
  }
>("deployments/cancel", async ({ deploymentId, withoutRollback }, thunkAPI) => {
  const { commandId } = await deploymentsApi.cancelDeployment({
    deploymentId,
    withoutRollback,
  });

  await thunkAPI.dispatch(fetchCommand(commandId));
});

export const deploymentsSlice = createSlice({
  name: "deployments",
  initialState: deploymentsAdapter.getInitialState<{
    loading: Record<string, boolean>;
    canceling: Record<string, boolean>;
    loadingList: boolean;
  }>({
    loading: {},
    canceling: {},
    loadingList: false,
  }),
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
        state.loadingList = true;
      })
      .addCase(fetchDeployments.fulfilled, (state, action) => {
        deploymentsAdapter.removeAll(state);
        if (action.payload.length > 0) {
          deploymentsAdapter.upsertMany(state, action.payload);
        }
        state.loadingList = false;
      })
      .addCase(fetchDeployments.rejected, (state) => {
        state.loadingList = false;
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
