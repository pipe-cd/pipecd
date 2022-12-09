import chartColor from "@material-ui/core/colors/blue";
import { FC } from "react";
import { InsightDataPoint, InsightStep } from "~/modules/insight";
import { ChartBase } from "../chart-base";

export interface DeploymentFrequencyChartProps {
  step: InsightStep;
  data: { name: string; points: InsightDataPoint.AsObject[] }[];
}

export const DeploymentFrequencyChart: FC<DeploymentFrequencyChartProps> = ({
  step,
  data,
}) => {
  return (
    <ChartBase
      title="Deployment Frequency"
      step={step}
      data={data}
      xName=""
      yName="Number of Deployments"
      yMax={undefined}
      lineColor={chartColor[500]}
      areaColor={chartColor[200]}
    />
  );
};
