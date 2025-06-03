import { FC } from "react";
import { ChartBase } from "../chart-base";
import { blue as chartColor } from "@mui/material/colors";
import { InsightDataPoint, InsightResolution } from "~~/model/insight_pb";

export interface DeploymentFrequencyChartProps {
  resolution: InsightResolution;
  data: { name: string; points: InsightDataPoint.AsObject[] }[];
}

export const DeploymentFrequencyChart: FC<DeploymentFrequencyChartProps> = ({
  resolution,
  data,
}) => {
  return (
    <ChartBase
      title="Deployment Frequency"
      resolution={resolution}
      data={data}
      xName=""
      yName="Number of Deployments"
      yMax={undefined}
      lineColor={chartColor[500]}
      areaColor={chartColor[200]}
    />
  );
};
