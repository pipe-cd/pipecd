import React, { FC, useCallback } from "react";
import { Box, makeStyles, Paper, Typography } from "@material-ui/core";
import { InsightDataPoint } from "pipe/pkg/app/web/model/insight_pb";
import {
  LineChart,
  Line,
  YAxis,
  XAxis,
  CartesianGrid,
  Tooltip,
} from "recharts";
import dayjs from "dayjs";
import { theme } from "../theme";

const useStyles = makeStyles((theme) => ({
  root: {
    padding: theme.spacing(3),
    width: 680,
  },
  chart: {
    flex: 1,
  },
}));

interface Props {
  data: InsightDataPoint.AsObject[];
}

const CHART_WIDTH = 600;
const CHART_HEIGHT = 400;

export const DeploymentFrequencyChart: FC<Props> = ({ data }) => {
  const classes = useStyles();

  const formatter = useCallback(
    (time: number | string) => dayjs(time).format("MMM DD"),
    []
  );

  return (
    <Paper elevation={1} className={classes.root}>
      <Typography variant="h6" component="div">
        Deployment Frequency
      </Typography>
      {data.length === 0 ? (
        <Box
          width={CHART_WIDTH}
          height={CHART_HEIGHT}
          display="flex"
          alignItems="center"
          justifyContent="center"
        >
          <Typography variant="body1">No data</Typography>
        </Box>
      ) : (
        <LineChart
          data={data}
          width={CHART_WIDTH}
          height={CHART_HEIGHT}
          margin={{
            top: 48,
            right: 24,
            left: 8,
            bottom: 8,
          }}
        >
          <CartesianGrid vertical={false} />
          <Line
            dataKey="value"
            stroke={theme.palette.primary.main}
            strokeWidth={2}
            isAnimationActive={false}
            dot={{ fill: theme.palette.primary.main }}
          />
          <YAxis />
          <XAxis dataKey="timestamp" tickFormatter={formatter} />
          <Tooltip labelFormatter={formatter} />
        </LineChart>
      )}
    </Paper>
  );
};
