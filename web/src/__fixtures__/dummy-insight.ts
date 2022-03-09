import dayjs from "dayjs";
import { InsightDataPoint } from "~/modules/insight";
import { randomNumber } from "./utils";

const today = dayjs();
export const dummyDataPointsList: InsightDataPoint.AsObject[] = Array.from(
  new Array(14)
)
  .map((_, i) => ({
    value: randomNumber(30),
    timestamp: today.subtract(i, "day").valueOf(),
  }))
  .reverse();

export function createInsightDataPointFromObject(
  o: InsightDataPoint.AsObject
): InsightDataPoint {
  const point = new InsightDataPoint();
  point.setTimestamp(o.timestamp);
  point.setValue(o.value);
  return point;
}

export function createDataPointsListFromObject(
  list: InsightDataPoint.AsObject[]
): InsightDataPoint[] {
  return list.map(createInsightDataPointFromObject);
}
