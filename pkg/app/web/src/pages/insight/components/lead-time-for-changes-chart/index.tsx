import chartColor from "@material-ui/core/colors/green";
import { FC } from "react";
import { InsightDataPoint } from "~/modules/insight";
import { ChartBase } from "../chart-base";

export interface LeadTimeForChangesChartProps {
  data: { name: string; points: InsightDataPoint.AsObject[] }[];
}

export const LeadTimeForChangesChart: FC<LeadTimeForChangesChartProps> = ({
  data,
}) => {
  return (
    <ChartBase
      title="Lead Time For Changes"
      data={data}
      yName="Lead time (min)"
      xName="Dates"
      lineColor={chartColor[500]}
      areaColor={chartColor[200]}
    />
  );
};
