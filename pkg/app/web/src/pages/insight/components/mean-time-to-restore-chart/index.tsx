import chartColor from "@material-ui/core/colors/yellow";
import { FC } from "react";
import { InsightDataPoint } from "~/modules/insight";
import { ChartBase } from "../chart-base";

export interface MeanTimeToRestoreChartProps {
  data: { name: string; points: InsightDataPoint.AsObject[] }[];
}

export const MeanTimeToRestoreChart: FC<MeanTimeToRestoreChartProps> = ({
  data,
}) => {
  return (
    <ChartBase
      title="Mean Time To Restore"
      data={data}
      yName="Restore time (min)"
      xName="Dates"
      lineColor={chartColor[500]}
      areaColor={chartColor[200]}
    />
  );
};
