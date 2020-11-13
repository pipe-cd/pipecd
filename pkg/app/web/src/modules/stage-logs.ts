import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import { LogBlock as LogBlockModel } from "pipe/pkg/app/web/model/logblock_pb";
import { AppState } from ".";
import { DeploymentStatus } from "../../../../../bazel-bin/pkg/app/web/model/deployment_pb";
import { getStageLog } from "../api/stage-log";
import { selectById as selectDeploymentById } from "./deployments";

export { LogSeverity } from "pipe/pkg/app/web/model/logblock_pb";

export type LogBlock = LogBlockModel.AsObject;

export type StageLog = {
  deploymentId: string;
  stageId: string;
  logBlocks: LogBlock[];
  completed: boolean;
};

// NOTE: Use deploymentId + stageId as record key.
type StageLogs = Record<string, StageLog>;
const initialState: StageLogs = {};

export const createActiveStageKey = (props: {
  deploymentId: string;
  stageId: string;
}): string => `${props.deploymentId}${props.stageId}`;

export const fetchStageLog = createAsyncThunk<
  StageLog,
  {
    deploymentId: string;
    stageId: string;
    offsetIndex: number;
    retriedCount: number;
  },
  { state: AppState }
>(
  "stage-logs/fetch",
  async ({ deploymentId, offsetIndex, retriedCount, stageId }, thunkAPI) => {
    const s = thunkAPI.getState();
    const deployment = selectDeploymentById(s.deployments, deploymentId);

    if (!deployment) {
      throw new Error(`Deployment: ${deploymentId} is not exists in state.`);
    }

    // When the Deployment Status is `Pending` and `Planned`, the log doesn't exist, so it returns an empty log instead of requesting it.
    if (
      deployment.status === DeploymentStatus.DEPLOYMENT_PLANNED ||
      deployment.status === DeploymentStatus.DEPLOYMENT_PENDING
    ) {
      return {
        stageId,
        deploymentId,
        logBlocks: [],
        completed: false,
      };
    }

    const response = await getStageLog({
      deploymentId,
      offsetIndex,
      retriedCount,
      stageId,
    });

    return {
      stageId,
      deploymentId,
      logBlocks: response.blocksList,
      completed: response.completed,
    };
  }
);

export const stageLogsSlice = createSlice({
  name: "stageLogs",
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchStageLog.pending, (state, action) => {
        const id = createActiveStageKey(action.meta.arg);
        if (state[id]) {
          state[id].completed = false;
        } else {
          state[id] = {
            stageId: action.meta.arg.stageId,
            deploymentId: action.meta.arg.deploymentId,
            logBlocks: [],
            completed: false,
          };
        }
      })
      .addCase(fetchStageLog.fulfilled, (state, action) => {
        const id = createActiveStageKey(action.meta.arg);
        state[id] = action.payload;
        state[id].completed = true;
      })
      .addCase(fetchStageLog.rejected, (state, action) => {
        const id = createActiveStageKey(action.meta.arg);
        state[id].completed = true;
      });
  },
});

export const selectStageLogById = (
  state: StageLogs,
  props: {
    deploymentId: string;
    stageId: string;
  }
): StageLog | null => {
  return state[createActiveStageKey(props)];
};
