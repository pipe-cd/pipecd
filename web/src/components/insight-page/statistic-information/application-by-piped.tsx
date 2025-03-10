import {
  Box,
  Card,
  CardContent,
  makeStyles,
  Typography,
} from "@material-ui/core";
import clsx from "clsx";
import { FC, useEffect, useMemo } from "react";
import {
  GridComponent,
  LegendComponent,
  TitleComponent,
  TooltipComponent,
} from "echarts/components";
import * as echarts from "echarts/core";
import { BarChart } from "echarts/charts";
import { CanvasRenderer } from "echarts/renderers";
import grey from "@material-ui/core/colors/grey";
import lineColor from "@material-ui/core/colors/purple";
import { useAppSelector } from "~/hooks/redux";
import { selectAll as selectAllApplications } from "~/modules/applications";
import { selectAllPipeds } from "~/modules/pipeds";
import ChartEmptyData from "~/components/chart-empty-data";
import useEChartState from "~/hooks/useEChartState";

const useStyles = makeStyles(() => ({
  root: {
    minWidth: 300,
    display: "inline-block",
    overflow: "visible",
    position: "relative",
  },
  pageTitle: {
    fontWeight: "bold",
  },
}));

const ApplicationByPiped: FC = () => {
  const classes = useStyles();
  const { chart, chartElm } = useEChartState({
    extensions: [
      TitleComponent,
      TooltipComponent,
      GridComponent,
      BarChart,
      CanvasRenderer,
      LegendComponent,
    ],
  });
  const applications = useAppSelector((state) =>
    selectAllApplications(state.applications)
  );

  const pipedList = useAppSelector(selectAllPipeds);

  const data: { name: string; count: number; rank: number }[] = useMemo(() => {
    const pipedMap = pipedList.reduce((acc, piped) => {
      acc[piped.id] = { name: piped.name, id: piped.id, count: 0 };
      return acc;
    }, {} as Record<string, { name: string; id: string; count: number }>);

    applications.forEach((app) => {
      if (!app.pipedId) return;

      if (app.pipedId in pipedMap === false) {
        return;
      }

      if (app.pipedId in pipedMap) {
        pipedMap[app.pipedId].count += 1;
      }
    });

    const listAppsByPipedSorted = Object.values(pipedMap)
      .filter((v) => v.count > 0)
      .sort((a, b) => b.count - a.count);

    const itemRank1st = listAppsByPipedSorted[0];
    const itemRank2nd = listAppsByPipedSorted[1];
    const itemRank3rd = listAppsByPipedSorted[2];

    const list = [];

    if (itemRank2nd) list.push({ ...itemRank2nd, rank: 2 });
    if (itemRank1st) list.push({ ...itemRank1st, rank: 1 });
    if (itemRank3rd) list.push({ ...itemRank3rd, rank: 3 });

    return list;
  }, [applications, pipedList]);

  const yMax = useMemo(() => {
    return Math.max(...data.map((v) => v.count));
  }, [data]);

  const isNoData = data.length === 0;

  useEffect(() => {
    if (chart && data.length !== 0) {
      chart.setOption({
        grid: {
          top: 50,
          bottom: 0,
          left: 0,
          right: 0,
        },
        xAxis: {
          data: data.map((v) => v.name),
          axisLine: { show: false },
          axisLabel: { show: false },
          axisTick: { show: false },
          splitLine: { show: false },
        },
        yAxis: {
          axisLine: { show: false },
          axisLabel: { show: false },
          axisTick: { show: false },
          splitLine: { show: false },
        },
        tooltip: { show: true },
        series: [
          {
            name: "Applications",
            type: "bar",
            stack: "title",
            data: data.map((v) => ({
              value: v.count,
              label: {
                show: true,
                formatter: v.rank === 1 ? "{b}\n{c} apps" : "{c} apps",
                align: "center",
                verticalAlign: "bottom",
                position: "top",
                distance: 10,
                width: 200,
                overflow: "truncate",
              },
              itemStyle: {
                color: v.rank === 1 ? lineColor[500] : grey[300],
                borderRadius: 5,
              },
            })),
          },
        ],
      } as echarts.EChartsInitOpts);
    }
  }, [chart, data, isNoData, yMax]);

  return (
    <Card raised className={clsx(classes.root)}>
      <CardContent>
        <Typography
          color="textSecondary"
          gutterBottom
          className={classes.pageTitle}
        >
          Application by piped
        </Typography>

        <Box position={"relative"}>
          <div style={{ width: "100%", height: 200 }} ref={chartElm} />
          <ChartEmptyData visible={!data.length} />
        </Box>
      </CardContent>
    </Card>
  );
};

export default ApplicationByPiped;
