import chartColor from "@material-ui/core/colors/blue";
import { InsightDataPoint } from "pipe/pkg/app/web/model/insight_pb";
import React, { FC } from "react";
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
      yName="Deployments"
      xName="Deployment Date"
      lineColor={chartColor[500]}
      areaColor={chartColor[200]}
    />
  );
};
