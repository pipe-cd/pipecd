import { Box, makeStyles, Typography } from "@material-ui/core";
import React, { FC, memo } from "react";

const useStyles = makeStyles((theme) => ({
  main: {
    flex: 1,
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
  },
  chart: {
    position: "absolute",
    display: "flex",
    alignItems: "flex-end",
  },
  bar: {
    width: 80,
    background: theme.palette.grey[300],
    margin: theme.spacing(1),
  },
  message: {
    zIndex: 10,
  },
}));

export const InsightIndexPage: FC = memo(function InsightIndexPage() {
  const classes = useStyles();
  return (
    <div className={classes.main}>
      <Typography variant="h2" className={classes.message}>
        COMING SOON
      </Typography>
      <div className={classes.chart}>
        <Box className={classes.bar} height={140} />
        <Box className={classes.bar} height={300} />
        <Box className={classes.bar} height={70} />
        <Box className={classes.bar} height={140} />
        <Box className={classes.bar} height={260} />
        <Box className={classes.bar} height={180} />
        <Box className={classes.bar} height={220} />
        <Box className={classes.bar} height={180} />
      </div>
    </div>
  );
});
