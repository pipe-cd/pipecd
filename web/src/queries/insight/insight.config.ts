import { InsightResolution } from "~~/model/insight_pb";
export {
  InsightMetricsKind,
  InsightDataPoint,
  InsightResolution,
} from "pipecd/web/model/insight_pb";

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

export const InsightResolutions = [
  InsightResolution.DAILY,
  InsightResolution.MONTHLY,
];

export const INSIGHT_RESOLUTION_TEXT: Record<InsightResolution, string> = {
  [InsightResolution.DAILY]: "Daily",
  [InsightResolution.MONTHLY]: "Monthly",
};
