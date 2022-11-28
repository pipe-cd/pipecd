import chartColor from "@material-ui/core/colors/blue";
import { FC } from "react";
import { InsightDataPoint } from "~/modules/insight";
import { ChartBase } from "../chart-base";

export interface DeploymentFrequencyChartProps {
  data: { name: string; points: InsightDataPoint.AsObject[] }[];
}

export const DeploymentFrequencyChart: FC<DeploymentFrequencyChartProps> = ({
  data,
}) => {
  return (
    <ChartBase
      title="Deployment Frequency"
      data={data}
      yName="Number of Deployments"
      xName="Deployment Date"
      lineColor={chartColor[500]}
      areaColor={chartColor[200]}
    />
  );
};
