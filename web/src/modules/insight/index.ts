import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import dayjs from "dayjs";

const MODULE_NAME = "insight";

export interface InsightState {
  rangeFrom: number;
  rangeTo: number;
  applicationId: string;
  // Suppose to be like ["key-1:value-1"]
  // sindresorhus/query-string doesn't support multidimensional arrays, that's why the format is a bit tricky.
  labels: Array<string>;
}

const now = dayjs(Date.now());

const initialState: InsightState = {
  rangeFrom: now.subtract(1, "month").valueOf(),
  rangeTo: now.valueOf(), //-new Date().getTimezoneOffset() * 60,
  applicationId: "",
  labels: [],
};

export const insightSlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {
    changeApplication(state, action: PayloadAction<string>) {
      state.applicationId = action.payload;
    },
    changeRangeFrom(state, action: PayloadAction<number>) {
      state.rangeFrom = action.payload;
    },
    changeRangeTo(state, action: PayloadAction<number>) {
      state.rangeTo = action.payload;
    },
  },
});

export const {
  changeApplication,
  changeRangeFrom,
  changeRangeTo,
} = insightSlice.actions;

export {
  InsightMetricsKind,
  InsightDataPoint,
} from "pipecd/web/model/insight_pb";
