import { Card, CardContent, Typography, Box } from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";
import {
  GridComponent,
  LegendComponent,
  TitleComponent,
  TooltipComponent,
} from "echarts/components";
import * as echarts from "echarts/core";
import clsx from "clsx";
import { FC, useEffect, useMemo } from "react";
import { useAppSelector } from "~/hooks/redux";
import { GaugeChart } from "echarts/charts";
import { CanvasRenderer } from "echarts/renderers";
import useEChartState from "~/hooks/useEChartState";
import dayjs from "dayjs";
import ChartEmptyData from "~/components/chart-empty-data";
import LegendRow from "./legend-row";
import { red } from "@mui/material/colors";

const failColor = red[500];

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

const Deployment24h: FC = () => {
  const classes = useStyles();
  const { chart, chartElm } = useEChartState({
    extensions: [
      TitleComponent,
      TooltipComponent,
      GridComponent,
      GaugeChart,
      CanvasRenderer,
      LegendComponent,
    ],
  });

  const deploymentFailRate = useAppSelector(
    (state) => state.deploymentChangeFailureRate.data24h
  );
  const deployment24h = useAppSelector(
    (state) => state.deploymentFrequency.data24h
  );

  const deploymentSummary = useMemo(() => {
    const summary = {
      totalDeployment: deployment24h?.[0]?.value || 0,
      failRate: (deploymentFailRate?.[0]?.value || 0) * 100,
      date: deploymentFailRate?.[0]?.timestamp
        ? dayjs.utc(deploymentFailRate[0].timestamp).format("DD/MM/YYYY")
        : "-",
    };
    return summary;
  }, [deployment24h, deploymentFailRate]);

  const data = useMemo(() => {
    return [
      {
        name: "Failed Rate",
        value: deploymentSummary.failRate,
        color: failColor,
        tooltip: { valueFormatter: (v: number) => v + "%" },
      },
    ];
  }, [deploymentSummary.failRate]);

  useEffect(() => {
    if (chart && data.length !== 0) {
      chart.setOption({
        color: [failColor],
        grid: {
          top: 0,
          bottom: 0,
          left: 0,
          right: 0,
        },
        title: {
          text: deploymentSummary?.totalDeployment?.toString(),
          left: "center",
          top: "center",
          textStyle: { fontSize: 30 },
          subtext: "Deployments",
          subtextGap: 5,
          itemGap: -10,
        },
        tooltip: { show: true },
        series: [
          {
            type: "gauge",
            startAngle: 90,
            endAngle: -270,
            radius: "100%",
            pointer: { show: false },
            progress: {
              show: true,
              overlap: false,
              roundCap: true,
              width: 30,
              clip: false,
            },
            axisLine: { lineStyle: { width: 20 } },
            splitLine: { show: false, distance: 0, length: 10 },
            axisTick: { show: false },
            axisLabel: { show: false, distance: 50 },
            data: data,
            title: { show: false },
            detail: { show: false },
          },
        ],
      } as echarts.EChartsCoreOption);
    }
  }, [chart, deploymentSummary.totalDeployment, data]);

  return (
    <Card raised className={clsx(classes.root)}>
      <CardContent>
        <Typography color="textSecondary" className={classes.pageTitle}>
          Deployments in 24h
        </Typography>
        <Box position={"relative"}>
          <div style={{ width: "100%", height: 150 }} ref={chartElm} />
          <ChartEmptyData visible={!deploymentSummary.totalDeployment} />
        </Box>

        <LegendRow
          data={[
            {
              key: "failRate",
              color: failColor,
              title: `${deploymentSummary.failRate}% failed`,
              description: `/ ${deploymentSummary.totalDeployment} deployments`,
            },
          ]}
        />
      </CardContent>
    </Card>
  );
};

export default Deployment24h;
