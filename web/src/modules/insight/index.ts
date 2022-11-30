import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { InsightStep } from "pipecd/web/model/insight_pb";
import dayjs from "dayjs";

const MODULE_NAME = "insight";

export enum InsightRange {
  LAST_1_WEEK = 0,
  LAST_1_MONTH = 1,
  LAST_3_MONTHS = 2,
  LAST_6_MONTHS = 3,
  LAST_1_YEAR = 4,
  LAST_2_YEARS = 5,
}

export const INSIGHT_RANGE_TEXT: Record<InsightRange, string> = {
  [InsightRange.LAST_1_WEEK]: "Last 1 week",
  [InsightRange.LAST_1_MONTH]: "Last 1 month",
  [InsightRange.LAST_3_MONTHS]: "Last 3 months",
  [InsightRange.LAST_6_MONTHS]: "Last 6 months",
  [InsightRange.LAST_1_YEAR]: "Last 1 year",
  [InsightRange.LAST_2_YEARS]: "Last 2 years",
};

export const InsightRanges = [
  InsightRange.LAST_1_WEEK,
  InsightRange.LAST_1_MONTH,
  InsightRange.LAST_3_MONTHS,
  InsightRange.LAST_6_MONTHS,
  InsightRange.LAST_1_YEAR,
  InsightRange.LAST_2_YEARS,
];

export const InsightSteps = [InsightStep.DAILY, InsightStep.MONTHLY];

export const INSIGHT_STEP_TEXT: Record<InsightStep, string> = {
  [InsightStep.DAILY]: "Daily",
  [InsightStep.MONTHLY]: "Monthly",
};

export interface InsightState {
  range: InsightRange;
  step: InsightStep;
  applicationId: string;
  // Suppose to be like ["key-1:value-1"]
  // sindresorhus/query-string doesn't support multidimensional arrays, that's why the format is a bit tricky.
  labels: Array<string>;
}

const initialState: InsightState = {
  range: InsightRange.LAST_1_MONTH,
  step: InsightStep.DAILY,
  applicationId: "",
  labels: [],
};

export const insightSlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {
    changeRange(state, action: PayloadAction<InsightRange>) {
      state.range = action.payload;
    },
    changeStep(state, action: PayloadAction<InsightStep>) {
      state.step = action.payload;
    },
    changeApplication(state, action: PayloadAction<string>) {
      state.applicationId = action.payload;
    },
    changeLabels(state, action: PayloadAction<Array<string>>) {
      state.labels = action.payload;
    },
  },
});

export const {
  changeRange,
  changeStep,
  changeApplication,
  changeLabels,
} = insightSlice.actions;

export {
  InsightMetricsKind,
  InsightDataPoint,
  InsightStep,
} from "pipecd/web/model/insight_pb";

export function makeTimeRange(r: InsightRange): [number, number] {
  const now = dayjs(Date.now());
  const rangeTo = now.valueOf();
  let rangeFrom = rangeTo;

  switch (r) {
    case InsightRange.LAST_1_WEEK:
      rangeFrom = now.subtract(1, "week").valueOf();
      break;
    case InsightRange.LAST_1_MONTH:
      rangeFrom = now.subtract(1, "month").valueOf();
      break;
    case InsightRange.LAST_3_MONTHS:
      rangeFrom = now.subtract(3, "month").valueOf();
      break;
    case InsightRange.LAST_6_MONTHS:
      rangeFrom = now.subtract(6, "month").valueOf();
      break;
    case InsightRange.LAST_1_YEAR:
      rangeFrom = now.subtract(1, "year").valueOf();
      break;
    case InsightRange.LAST_2_YEARS:
      rangeFrom = now.subtract(2, "year").valueOf();
      break;
  }

  return [rangeFrom, rangeTo];
}
