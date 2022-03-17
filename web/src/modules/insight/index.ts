import { createSlice, PayloadAction } from "@reduxjs/toolkit";
// import { InsightStep } from "pipecd/web/model/insight_pb";
import dayjs from "dayjs";

const MODULE_NAME = "insight";

export interface InsightState {
  applicationId: string;
  // step: InsightStep;
  rangeFrom: number;
  rangeTo: number;
}

const now = dayjs(Date.now());

const initialState: InsightState = {
  applicationId: "",
  // step: InsightStep.DAILY,
  rangeFrom: now.valueOf(),
  rangeTo: now.add(7, "day").valueOf(),
};

export const insightSlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {
    changeApplication(state, action: PayloadAction<string>) {
      state.applicationId = action.payload;
    },
    // changeStep(state, action: PayloadAction<InsightStep>) {
    //   state.step = action.payload;
    // },
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
  // changeStep,
  changeRangeFrom,
  changeRangeTo,
} = insightSlice.actions;

export {
  InsightMetricsKind,
  // InsightStep,
  InsightDataPoint,
} from "pipecd/web/model/insight_pb";
