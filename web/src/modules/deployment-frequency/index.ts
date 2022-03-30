import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import type { AppState } from "~/store";
import { InsightMetricsKind, InsightDataPoint } from "../insight";
import * as InsightAPI from "~/api/insight";
import { LoadingStatus } from "~/types/module";
import { InsightResultType } from "pipecd/web/model/insight_pb";

const MODULE_NAME = "deploymentFrequency";

export interface DeploymentFrequencyState {
  status: LoadingStatus;
  data: InsightDataPoint.AsObject[];
}

const initialState: DeploymentFrequencyState = {
  status: "idle",
  data: [],
};

export const fetchDeploymentFrequency = createAsyncThunk<
  InsightDataPoint.AsObject[],
  void,
  { state: AppState }
>(`${MODULE_NAME}/fetch`, async (_, thunkAPI) => {
  const state = thunkAPI.getState();

  const data = await InsightAPI.getInsightData({
    applicationId: state.insight.applicationId,
    metricsKind: InsightMetricsKind.DEPLOYMENT_FREQUENCY,
    rangeFrom: state.insight.rangeFrom,
    rangeTo: state.insight.rangeTo,
    timezone: state.insight.timezone,
  });

  if (data.type == InsightResultType.MATRIX) {
    return data.matrixList[0].dataPointsList;
  } else {
    return [];
  }
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
