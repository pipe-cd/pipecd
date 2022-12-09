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
  // Find the best yMax value to make the graph more readable.
  let yMax = -1;
  data.forEach((d) => {
    d.points.forEach((p) => {
      if (p.value > yMax) {
        yMax = p.value;
      }
    });
  });
  if (yMax > 0.1) {
    yMax = 1.0;
  } else if (yMax > 0.05) {
    yMax = 0.5;
  } else {
    yMax = 0.1;
  }

  return (
    <ChartBase
      title="Change Failure Rate"
      step={step}
      data={data}
      xName=""
      yName="Failed Deployments / Total"
      yMax={yMax}
      lineColor={chartColor[500]}
      areaColor={chartColor[200]}
    />
  );
};
