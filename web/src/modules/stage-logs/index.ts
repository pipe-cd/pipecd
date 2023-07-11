import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import { LogBlock } from "pipecd/web/model/logblock_pb";
import { getStageLog } from "~/api/stage-log";
import {
  selectById as selectDeploymentById,
  StageStatus,
} from "../deployments";
import { AppState } from "~/store";
import { StatusCode } from "grpc-web";

export type StageLog = {
  deploymentId: string;
  stageId: string;
  logBlocks: LogBlock.AsObject[];
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
    const initialLogs: StageLog = {
      stageId,
      deploymentId,
      logBlocks: [],
    };

    if (!deployment) {
      throw new Error(`Deployment: ${deploymentId} is not exists in state.`);
    }

    const stage = deployment.stagesList.find((stage) => stage.id === stageId);
    if (!stage) {
      throw new Error(
        `Stage (ID: ${stageId}) is not found in application state.`
      );
    }

    if (stage.status === StageStatus.STAGE_NOT_STARTED_YET) {
      return initialLogs;
    }

    const response = await getStageLog({
      deploymentId,
      offsetIndex,
      retriedCount,
      stageId,
    }).catch((e: { code: number }) => {
      // If status is running and error code is NOT_FOUND, it is maybe first state of deployment log.
      // So we ignore this error and then return initialLogs below code.
      if (
        e.code === StatusCode.NOT_FOUND &&
        stage.status === StageStatus.STAGE_RUNNING
      ) {
        return;
      }

      throw e;
    });

    if (!response) {
      return initialLogs;
    }

    return {
      stageId,
      deploymentId,
      logBlocks: response.blocksList,
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
        if (!state[id]) {
          state[id] = {
            stageId: action.meta.arg.stageId,
            deploymentId: action.meta.arg.deploymentId,
            logBlocks: [],
          };
        }
      })
      .addCase(fetchStageLog.fulfilled, (state, action) => {
        const id = createActiveStageKey(action.meta.arg);
        // Skip update state when no log updates
        // to avoid unnecessary state update and trigger re-render.
        if (JSON.stringify(state[id]) !== JSON.stringify(action.payload)) {
          state[id] = action.payload;
        }
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

export { LogBlock, LogSeverity } from "pipecd/web/model/logblock_pb";
