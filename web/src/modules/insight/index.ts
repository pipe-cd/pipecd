import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import dayjs from "dayjs";

const MODULE_NAME = "insight";

export interface InsightState {
  applicationId: string;
  rangeFrom: number;
  rangeTo: number;
  timezone: string;
}

const now = dayjs(Date.now());

const initialState: InsightState = {
  applicationId: "",
  rangeFrom: now.subtract(1, "month").valueOf(),
  rangeTo: now.valueOf(),
  timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
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
