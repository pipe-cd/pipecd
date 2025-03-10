import {
  Box,
  Card,
  CardContent,
  makeStyles,
  Typography,
} from "@material-ui/core";
import clsx from "clsx";
import { FC, useEffect, useMemo } from "react";
import { useAppSelector } from "~/hooks/redux";
import red from "@material-ui/core/colors/red";
import green from "@material-ui/core/colors/green";
import useEChartState from "~/hooks/useEChartState";
import { PieChart } from "echarts/charts";
import {
  GridComponent,
  LegendComponent,
  TitleComponent,
  TooltipComponent,
} from "echarts/components";
import { CanvasRenderer } from "echarts/renderers";
import ChartEmptyData from "~/components/chart-empty-data";
import LegendRow from "./legend-row";
import { getPercentage } from "~/utils/common";

const enabledColor = green[500];
const disabledColor = red[500];

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

const ApplicationCount: FC = () => {
  const classes = useStyles();
  const appSummary = useAppSelector((state) => state.applicationCounts.summary);

  const { chart, chartElm } = useEChartState({
    extensions: [
      TitleComponent,
      TooltipComponent,
      GridComponent,
      PieChart,
      CanvasRenderer,
      LegendComponent,
    ],
  });

  const data = useMemo(() => {
    return [
      {
        name: "Enabled apps",
        value: getPercentage(appSummary.enabled, appSummary.total, 2),
        progress: { itemStyle: { color: enabledColor } },
        color: enabledColor,
        tooltip: { valueFormatter: (v: number) => v + "%" },
      },
      {
        name: "Disabled apps",
        value: getPercentage(appSummary.disabled, appSummary.total, 2),
        progress: { itemStyle: { color: disabledColor } },
        color: enabledColor,
        tooltip: { valueFormatter: (v: number) => v + "%" },
      },
    ];
  }, [appSummary.total, appSummary.disabled, appSummary.enabled]);

  useEffect(() => {
    if (chart && data.length !== 0) {
      chart.setOption({
        color: [enabledColor, disabledColor],
        grid: { top: 0, bottom: 0, left: 0, right: 0 },
        title: {
          text: appSummary.total.toString(),
          left: "center",
          top: "center",
          textStyle: { fontSize: 30 },
          subtext: "apps",
          subtextGap: 5,
          itemGap: -10,
        },
        tooltip: { trigger: "item" },
        series: [
          {
            name: "Applications",
            type: "pie",
            radius: ["70%", "100%"],
            itemStyle: {
              borderRadius: 10,
              borderColor: "#fff",
              borderWidth: 2,
            },
            label: { show: false },
            data,
          },
        ],
      });
    }
  }, [chart, data, appSummary.total]);

  return (
    <Card raised className={clsx(classes.root)}>
      <CardContent>
        <Typography color="textSecondary" className={classes.pageTitle}>
          Applications
        </Typography>
        <Box position={"relative"}>
          <div style={{ width: "100%", height: 150 }} ref={chartElm} />
          <ChartEmptyData visible={!appSummary.total} />
        </Box>

        <LegendRow
          data={[
            {
              key: "enabled",
              color: enabledColor,
              title: `${appSummary.enabled} enabled`,
              description: `/ ${appSummary.total} apps`,
            },
            {
              key: "disabled",
              color: disabledColor,
              title: `${appSummary.disabled} disabled`,
              description: `/ ${appSummary.total} apps`,
            },
          ]}
        />
      </CardContent>
    </Card>
  );
};

export default ApplicationCount;
