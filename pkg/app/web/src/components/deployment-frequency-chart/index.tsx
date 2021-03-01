import { Box, makeStyles, Paper, Typography } from "@material-ui/core";
import { WarningOutlined } from "@material-ui/icons";
import dayjs from "dayjs";
import { LineChart } from "echarts/charts";
import {
  GridComponent,
  LegendComponent,
  TitleComponent,
  TooltipComponent,
} from "echarts/components";
import * as echarts from "echarts/core";
import { CanvasRenderer } from "echarts/renderers";
import { InsightDataPoint } from "pipe/pkg/app/web/model/insight_pb";
import React, { FC, useEffect, useState } from "react";

echarts.use([
  TitleComponent,
  TooltipComponent,
  GridComponent,
  LineChart,
  CanvasRenderer,
  LegendComponent,
]);

const useStyles = makeStyles((theme) => ({
  root: {
    minWidth: 600,
  },
  noDataMessage: {
    display: "flex",
  },
  noDataMessageIcon: {
    marginRight: theme.spacing(1),
  },
  title: {
    padding: theme.spacing(3),
    paddingBottom: 0,
  },
}));

export interface DeploymentFrequencyChartProps {
  data: { name: string; points: InsightDataPoint.AsObject[] }[];
}

const labelFormatter = (time: number | string): string =>
  dayjs(time).format("YYYY MMM DD");

const NO_DATA_TEXT = "No data is available.";

export const DeploymentFrequencyChart: FC<DeploymentFrequencyChartProps> = ({
  data,
}) => {
  const classes = useStyles();
  const [chart, setChart] = useState<echarts.ECharts | null>(null);

  useEffect(() => {
    if (chart && data.length !== 0) {
      chart.setOption({
        legend: { data: data.map((v) => v.name) },
        xAxis: [
          {
            type: "category",
            boundaryGap: false,
            data: data[0].points.map((data) => labelFormatter(data.timestamp)),
          },
        ],
        tooltip: {
          trigger: "axis",
        },
        yAxis: [{ type: "value" }],
        series: data.map((v) => ({
          name: v.name,
          type: "line",
          stack: "deployments",
          data: v.points.map((point) => point.value),
          areaStyle: {},
          emphasis: {
            focus: "series",
          },
        })),
      });
    }
  }, [chart, data]);

  return (
    <Paper elevation={1} className={classes.root}>
      <Typography variant="h6" component="div" className={classes.title}>
        Deployment Frequency
      </Typography>
      {data.length === 0 ? (
        <Box
          display="flex"
          alignItems="center"
          justifyContent="center"
          height={420}
        >
          <Typography
            variant="body1"
            color="textSecondary"
            className={classes.noDataMessage}
          >
            <WarningOutlined className={classes.noDataMessageIcon} />
            {NO_DATA_TEXT}
          </Typography>
        </Box>
      ) : (
        <div
          style={{ width: "100%", height: 400 }}
          ref={(ref) => {
            if (ref) {
              setChart(echarts.init(ref));
            }
          }}
        />
      )}
    </Paper>
  );
};
