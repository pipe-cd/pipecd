import chartColor from "@material-ui/core/colors/deepPurple";
import { FC } from "react";
import { InsightDataPoint } from "~/modules/insight";
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
      yName="Failed Deployments / Total"
      xName="Release Date"
      lineColor={chartColor[500]}
      areaColor={chartColor[200]}
    />
  );
};
