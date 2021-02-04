import React, { FC } from "react";
import { Box, makeStyles, Paper, Typography } from "@material-ui/core";
import { InsightDataPoint } from "pipe/pkg/app/web/model/insight_pb";
import {
  LineChart,
  Line,
  YAxis,
  XAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import dayjs from "dayjs";
import { theme } from "../../theme";
import { WarningOutlined } from "@material-ui/icons";

const useStyles = makeStyles((theme) => ({
  root: {
    padding: theme.spacing(3),
    minWidth: 600,
  },
  noDataMessage: {
    display: "flex",
  },
  noDataMessageIcon: {
    marginRight: theme.spacing(1),
  },
}));

export interface DeploymentFrequencyChartProps {
  data: InsightDataPoint.AsObject[];
}

const tickFormatter = (time: number | string): string =>
  dayjs(time).format("MMM DD");

const labelFormatter = (time: number | string): string =>
  dayjs(time).format("YYYY MMM DD");

const NO_DATA_TEXT = "No data is available.";

export const DeploymentFrequencyChart: FC<DeploymentFrequencyChartProps> = ({
  data,
}) => {
  const classes = useStyles();

  return (
    <Paper elevation={1} className={classes.root}>
      <Typography variant="h6" component="div">
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
        <ResponsiveContainer width="100%" aspect={2}>
          <LineChart
            data={data}
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
            <YAxis axisLine={false} tickLine={false} />
            <XAxis dataKey="timestamp" tickFormatter={tickFormatter} />
            <Tooltip labelFormatter={labelFormatter} />
          </LineChart>
        </ResponsiveContainer>
      )}
    </Paper>
  );
};
