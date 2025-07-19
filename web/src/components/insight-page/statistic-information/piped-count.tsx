import { CardContent, Typography, Box } from "@mui/material";
import {
  GridComponent,
  LegendComponent,
  TitleComponent,
  TooltipComponent,
} from "echarts/components";
import { FC, useEffect, useMemo } from "react";
import { Piped } from "~/modules/pipeds";
import { GaugeChart } from "echarts/charts";
import { CanvasRenderer } from "echarts/renderers";
import useEChartState from "~/hooks/useEChartState";
import ChartEmptyData from "~/components/chart-empty-data";
import LegendRow from "./legend-row";

import { blue as cyan, green } from "@mui/material/colors";
import { CardWrapper } from "./styles";
import { useGetPipeds } from "~/queries/pipeds/use-get-pipeds";

const enabledColor = cyan[500];
const onlineColor = green[500];

const PipedCount: FC = () => {
  const { data: pipeds = [] } = useGetPipeds({ withStatus: true });

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

  const pipedSummary = useMemo(() => {
    let enabledCount = 0;
    let onlineCount = 0;
    let totalPipeds = 0;

    pipeds.forEach((element) => {
      totalPipeds += 1;
      if (!element.disabled) enabledCount += 1;
      if (element.status === Piped.ConnectionStatus.ONLINE) onlineCount += 1;
    });
    return {
      total: totalPipeds,
      enabled: enabledCount,
      online: onlineCount,
      enabledPercent: (enabledCount / (totalPipeds || 1)) * 100,
      onlinePercent: (onlineCount / (enabledCount || 1)) * 100,
    };
  }, [pipeds]);

  const gaugeData = useMemo(() => {
    const enabledTooltip = `${pipedSummary.enabled} / ${pipedSummary.total} pipeds`;
    const onlineTooltip = `${pipedSummary.online} / ${pipedSummary.enabled} pipeds`;
    return [
      {
        name: "Enabled",
        value: pipedSummary.enabledPercent,
        progress: { itemStyle: { color: enabledColor } },
        tooltip: { valueFormatter: () => enabledTooltip },
      },
      {
        name: "Online",
        value: pipedSummary.onlinePercent,
        progress: { itemStyle: { color: onlineColor } },
        tooltip: { valueFormatter: () => onlineTooltip },
      },
    ];
  }, [
    pipedSummary.enabled,
    pipedSummary.enabledPercent,
    pipedSummary.online,
    pipedSummary.onlinePercent,
    pipedSummary.total,
  ]);

  useEffect(() => {
    if (chart && gaugeData.length !== 0) {
      chart.setOption({
        color: [enabledColor, onlineColor],
        grid: { top: 0, bottom: 0, left: 0, right: 0 },
        title: {
          text: pipedSummary.total.toString(),
          left: "center",
          top: "center",
          textStyle: { fontSize: 30 },
          subtext: "Piped",
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
              clip: false,
            },
            axisLine: { lineStyle: { width: 40 } },
            splitLine: { show: false, distance: 0, length: 10 },
            axisTick: { show: false },
            axisLabel: { show: false, distance: 50 },
            data: gaugeData,
            title: { show: false },
            detail: { show: false },
          },
        ],
      });
    }
  }, [chart, gaugeData, pipedSummary.total]);

  return (
    <CardWrapper raised>
      <CardContent>
        <Typography
          color="textSecondary"
          sx={{
            fontWeight: "bold",
          }}
        >
          Piped
        </Typography>
        <Box
          sx={{
            position: "relative",
          }}
        >
          <div style={{ width: "100%", height: 150 }} ref={chartElm} />
          <ChartEmptyData visible={!pipedSummary.total} />
        </Box>

        <LegendRow
          data={[
            {
              key: "enabled",
              color: enabledColor,
              title: `${pipedSummary.enabled} enabled`,
              description: `/ ${pipedSummary.total} total`,
            },
            {
              key: "online",
              color: onlineColor,
              title: `${pipedSummary.online} online`,
              description: `/ ${pipedSummary.enabled} enabled`,
            },
          ]}
        />
      </CardContent>
    </CardWrapper>
  );
};

export default PipedCount;
