import chartColor from "@material-ui/core/colors/deepPurple";
import { FC } from "react";
import { InsightDataPoint, InsightStep } from "~/modules/insight";
import { ChartBase } from "../chart-base";

export interface ChangeFailureRateChartProps {
  step: InsightStep;
  data: { name: string; points: InsightDataPoint.AsObject[] }[];
}

export const ChangeFailureRateChart: FC<ChangeFailureRateChartProps> = ({
  step,
  data,
}) => {
  return (
    <ChartBase
      title="Change Failure Rate"
      step={step}
      data={data}
      yName="Failed Deployments / Total"
      xName=""
      lineColor={chartColor[500]}
      areaColor={chartColor[200]}
    />
  );
};
