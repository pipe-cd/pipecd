import dayjs from "dayjs";
import utc from "dayjs/plugin/utc";
import { InsightRange } from "~/queries/insight/insight.config";
import { InsightResolution } from "~~/model/insight_pb";

export function determineTimeRange(
  r: InsightRange,
  s: InsightResolution
): [number, number] {
  // Load utc plugin.
  dayjs.extend(utc);

  const rangeTo = dayjs.utc().endOf("day");
  let rangeFrom = rangeTo;

  if (s === InsightResolution.DAILY) {
    switch (r) {
      case InsightRange.LAST_1_WEEK:
        rangeFrom = rangeTo.subtract(7, "day");
        break;
      case InsightRange.LAST_1_MONTH:
        rangeFrom = rangeTo.subtract(1, "month");
        break;
      case InsightRange.LAST_3_MONTHS:
        rangeFrom = rangeTo.subtract(3, "month");
        break;
      case InsightRange.LAST_6_MONTHS:
        rangeFrom = rangeTo.subtract(6, "month");
        break;
      case InsightRange.LAST_1_YEAR:
        rangeFrom = rangeTo.subtract(1, "year");
        break;
      case InsightRange.LAST_2_YEARS:
        rangeFrom = rangeTo.subtract(2, "year");
        break;
    }
    rangeFrom = rangeFrom.add(1, "day").startOf("day");
  } else {
    switch (r) {
      case InsightRange.LAST_3_MONTHS:
        rangeFrom = rangeTo.subtract(2, "month");
        break;
      case InsightRange.LAST_6_MONTHS:
        rangeFrom = rangeTo.subtract(5, "month");
        break;
      case InsightRange.LAST_1_YEAR:
        rangeFrom = rangeTo.subtract(11, "month");
        break;
      case InsightRange.LAST_2_YEARS:
        rangeFrom = rangeTo.subtract(23, "month");
        break;
    }
    rangeFrom = rangeFrom.startOf("month");
  }

  return [rangeFrom.valueOf(), rangeTo.valueOf()];
}
