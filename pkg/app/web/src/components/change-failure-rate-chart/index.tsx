import chartColor from "@material-ui/core/colors/deepPurple";
import { InsightDataPoint } from "pipe/pkg/app/web/model/insight_pb";
import React, { FC } from "react";
import { ChartBase } from "../chart-base";

export interface ChangeFailureRateChartProps {
  data: { name: string; points: InsightDataPoint.AsObject[] }[];
}

export const ChangeFailureRateChart: FC<ChangeFailureRateChartProps> = ({
  data,
}) => {
  return (
    <ChartBase
      title="Change Failure Rate"
      data={data}
      yName="Failed release attempts (%)"
      xName="Release Date"
      lineColor={chartColor[500]}
      areaColor={chartColor[200]}
    />
  );
};
