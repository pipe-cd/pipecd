import { Box, makeStyles, Paper, Typography } from "@material-ui/core";
import grey from "@material-ui/core/colors/grey";
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
import { FC, useCallback, useEffect, useRef, useState } from "react";
import { InsightDataPoint } from "~/modules/insight";
import { dummyDataPointsList } from "~/__fixtures__/dummy-insight";

const placeholderData = [{ name: "All", points: dummyDataPointsList }];

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
    position: "relative",
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

const labelFormatter = (time: number | string): string =>
  dayjs(time).format("YYYY MMM DD");

const NO_DATA_TEXT = "No data is available.";

export interface ChartBaseProps {
  title: string;
  xName: string;
  yName: string;
  data: { name: string; points: InsightDataPoint.AsObject[] }[];
  lineColor: string;
  areaColor: string;
}

const tooltip = {
  trigger: "axis",
};

export const ChartBase: FC<ChartBaseProps> = ({
  title,
  xName,
  yName,
  data,
  lineColor,
  areaColor,
}) => {
  const classes = useStyles();
  const [chart, setChart] = useState<echarts.ECharts | null>(null);
  const chartElm = useRef<HTMLDivElement | null>(null);
  const isNoData = data.length === 0;
  const _data = isNoData ? placeholderData : data;

  useEffect(() => {
    if (chart && _data.length !== 0) {
      chart.setOption({
        legend: { data: _data.map((v) => v.name) },
        xAxis: {
          type: "category",
          name: xName,
          nameLocation: "center",
          nameGap: 32,
          boundaryGap: false,
          data: _data[0].points.map((v) => labelFormatter(v.timestamp)),
        },
        yAxis: {
          type: "value",
          name: yName,
          nameLocation: "center",
          nameGap: 32,
        },
        tooltip,
        series: _data.map((v) => ({
          name: v.name,
          type: "line",
          stack: title,
          data: v.points.map((point) => point.value),
          emphasis: {
            focus: "series",
          },
          itemStyle: {
            color: isNoData ? grey[300] : lineColor,
          },
          lineStyle: {
            color: isNoData ? grey[300] : lineColor,
          },
          areaStyle: {
            color: isNoData ? grey[300] : areaColor,
          },
        })),
      });
    }
  }, [chart, _data, lineColor, areaColor, xName, yName, title, isNoData]);

  useEffect(() => {
    if (chartElm.current) {
      setChart(echarts.init(chartElm.current));
    }
  }, [chartElm]);

  const handleResize = useCallback(() => {
    if (chart) {
      chart.resize();
    }
  }, [chart]);

  useEffect(() => {
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, [handleResize]);

  return (
    <Paper elevation={1} className={classes.root}>
      <Typography variant="h6" component="div" className={classes.title}>
        {title}
      </Typography>

      <div style={{ width: "100%", height: 400 }} ref={chartElm} />
      {data.length === 0 ? (
        <Box
          width="100%"
          height="100%"
          display="flex"
          alignItems="center"
          justifyContent="center"
          position="absolute"
          top={0}
          left={0}
          bgcolor="#fafafabb"
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
      ) : null}
    </Paper>
  );
};
