import { Box, CardContent, Typography } from "@mui/material";
import { FC, useEffect, useMemo } from "react";
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

import { red, green } from "@mui/material/colors";
import { CardWrapper } from "./styles";
import { useGetApplicationCounts } from "~/queries/application-counts/use-get-application-counts";

const enabledColor = green[500];
const disabledColor = red[500];

const ApplicationCount: FC = () => {
  const { data: queryData } = useGetApplicationCounts();

  const appSummary = useMemo(() => {
    return {
      total: queryData?.summary.total ?? 0,
      enabled: queryData?.summary.enabled ?? 0,
      disabled: queryData?.summary.disabled ?? 0,
    };
  }, [
    queryData?.summary.disabled,
    queryData?.summary.enabled,
    queryData?.summary.total,
  ]);

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
    <CardWrapper raised>
      <CardContent>
        <Typography
          color="textSecondary"
          sx={{
            fontWeight: "bold",
          }}
        >
          Applications
        </Typography>
        <Box
          sx={{
            position: "relative",
          }}
        >
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
    </CardWrapper>
  );
};

export default ApplicationCount;
