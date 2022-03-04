import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import type { AppState } from "~/store";
import { InsightMetricsKind, InsightDataPoint } from "../insight";
import dayjs from "dayjs";
import * as InsightAPI from "~/api/insight";
import { InsightStep } from "pipecd/pkg/app/web/model/insight_pb";
import { LoadingStatus } from "~/types/module";

const MODULE_NAME = "deploymentFrequency";

export interface DeploymentFrequencyState {
  status: LoadingStatus;
  data: InsightDataPoint.AsObject[];
}

const initialState: DeploymentFrequencyState = {
  status: "idle",
  data: [],
};

const STEP_UNIT_MAP: Record<InsightStep, "day" | "week" | "month" | "year"> = {
  [InsightStep.DAILY]: "day",
  [InsightStep.WEEKLY]: "week",
  [InsightStep.MONTHLY]: "month",
  [InsightStep.YEARLY]: "year",
};

export const fetchDeploymentFrequency = createAsyncThunk<
  InsightDataPoint.AsObject[],
  void,
  { state: AppState }
>(`${MODULE_NAME}/fetch`, async (_, thunkAPI) => {
  const state = thunkAPI.getState();

  const { dataPointsList } = await InsightAPI.getInsightData({
    applicationId: state.insight.applicationId,
    step: state.insight.step,
    dataPointCount: dayjs(state.insight.rangeTo).diff(
      state.insight.rangeFrom,
      STEP_UNIT_MAP[state.insight.step]
    ),
    metricsKind: InsightMetricsKind.DEPLOYMENT_FREQUENCY,
    rangeFrom: state.insight.rangeFrom,
  });
  return dataPointsList;
});

export const deploymentFrequencySlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchDeploymentFrequency.pending, (state) => {
        state.status = "loading";
      })
      .addCase(fetchDeploymentFrequency.rejected, (state) => {
        state.status = "failed";
      })
      .addCase(fetchDeploymentFrequency.fulfilled, (state, action) => {
        state.status = "succeeded";
        state.data = action.payload;
      });
  },
});
