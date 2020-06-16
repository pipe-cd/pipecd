import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import { LogBlock as LogBlockModel } from "pipe/pkg/app/web/model/logblock_pb";
import { getStageLog } from "../api/stage-log";

type LogBlock = LogBlockModel.AsObject;

export type StageLog = {
  deploymentId: string;
  stageId: string;
  logBlocks: LogBlock[];
  completed: boolean;
};

// NOTE: Use deploymentId + stageId as record key.
type StageLogs = Record<string, StageLog>;
const initialState: StageLogs = {};

export const fetchStageLog = createAsyncThunk<
  StageLog,
  {
    deploymentId: string;
    stageId: string;
    offsetIndex: number;
    retriedCount: number;
  }
>(
  "stage-logs/fetch",
  async ({ deploymentId, offsetIndex, retriedCount, stageId }) => {
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
        state[`${action.meta.arg.deploymentId}/${action.meta.arg.stageId}`] = {
          stageId: action.meta.arg.stageId,
          deploymentId: action.meta.arg.deploymentId,
          logBlocks: [],
          completed: false,
        };
      })
      .addCase(fetchStageLog.fulfilled, (state, action) => {
        state[`${action.meta.arg.deploymentId}/${action.meta.arg.stageId}`] =
          action.payload;
      })
      .addCase(fetchStageLog.rejected, (state, action) => {});
  },
});

export const selectStageLogById = (
  state: StageLogs,
  {
    deploymentId,
    offsetIndex,
  }: {
    deploymentId: string;
    offsetIndex: string;
  }
) => {
  return state[`${deploymentId}/${offsetIndex}`];
};
