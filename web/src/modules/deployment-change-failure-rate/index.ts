import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import type { AppState } from "~/store";
import { InsightMetricsKind, InsightDataPoint } from "../insight";
import * as InsightAPI from "~/api/insight";
import { LoadingStatus } from "~/types/module";
import { InsightResultType } from "pipecd/web/model/insight_pb";

const MODULE_NAME = "deploymentChangeFailureRate";

export interface DeploymentChangeFailureRateState {
  status: LoadingStatus;
  data: InsightDataPoint.AsObject[];
}

const initialState: DeploymentChangeFailureRateState = {
  status: "idle",
  data: [],
};

export const fetchDeploymentChangeFailureRate = createAsyncThunk<
  InsightDataPoint.AsObject[],
  void,
  { state: AppState }
>(`${MODULE_NAME}/fetch`, async (_, thunkAPI) => {
  const state = thunkAPI.getState();

  const labels = new Array<[string, string]>();
  if (state.insight.labels) {
    for (const label of state.insight.labels) {
      const pair = label.split(":");
      pair.length === 2 && labels.push([pair[0], pair[1]]);
    }
  }

  const data = await InsightAPI.getInsightData({
    metricsKind: InsightMetricsKind.CHANGE_FAILURE_RATE,
    rangeFrom: state.insight.rangeFrom,
    rangeTo: state.insight.rangeTo,
    applicationId: state.insight.applicationId,
    labelsMap: labels,
  });

  if (data.type == InsightResultType.MATRIX) {
    return data.matrixList[0].dataPointsList;
  } else {
    return [];
  }
});

export const deploymentChangeFailureRateSlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchDeploymentChangeFailureRate.pending, (state) => {
        state.status = "loading";
      })
      .addCase(fetchDeploymentChangeFailureRate.rejected, (state) => {
        state.status = "failed";
      })
      .addCase(fetchDeploymentChangeFailureRate.fulfilled, (state, action) => {
        state.status = "succeeded";
        state.data = action.payload;
      });
  },
});
