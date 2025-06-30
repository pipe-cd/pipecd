import { FC, useMemo } from "react";
import { ChartBase } from "../chart-base";
import { deepPurple as chartColor } from "@mui/material/colors";
import {
  InsightDataPoint,
  InsightResolution,
} from "~/queries/insight/insight.config";

export interface ChangeFailureRateChartProps {
  resolution: InsightResolution;
  data: { name: string; points: InsightDataPoint.AsObject[] }[];
}

export const ChangeFailureRateChart: FC<ChangeFailureRateChartProps> = ({
  resolution,
  data,
}) => {
  // Find the best yMax value to make the graph more readable.
  const yMax = useMemo(() => {
    let max = -1;
    data.forEach((d) => {
      d.points.forEach((p) => {
        if (p.value > max) {
          max = p.value;
        }
      });
    });
    if (max > 0.1) {
      max = 1.0;
    } else if (max > 0.05) {
      max = 0.5;
    } else {
      max = 0.1;
    }
    return max;
  }, [data]);

  return (
    <ChartBase
      title="Change Failure Rate"
      resolution={resolution}
      data={data}
      xName=""
      yName="Failed Deployments / Total"
      yMax={yMax}
      lineColor={chartColor[500]}
      areaColor={chartColor[200]}
    />
  );
};
