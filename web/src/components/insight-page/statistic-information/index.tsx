import { makeStyles } from "@material-ui/core";
import { FC, memo } from "react";
import ApplicationCount from "./application-count";
import ApplicationByPiped from "./application-by-piped";
import PipedCount from "./piped-count";
import Deployment24h from "./deployment-24h";
import clsx from "clsx";

const useStyles = makeStyles((theme) => ({
  root: {
    marginBottom: theme.spacing(2),
  },
  group: {
    display: "flex",
    justifyContent: "center",
    gap: theme.spacing(2),
    flexWrap: "wrap",
  },
}));

export const StatisticInformation: FC = memo(function StatisticInformation() {
  const classes = useStyles();

  return (
    <div className={clsx(classes.root, classes.group)}>
      <div className={classes.group}>
        <ApplicationCount />
        <ApplicationByPiped />
      </div>
      <div className={classes.group}>
        <PipedCount />
        <Deployment24h />
      </div>
    </div>
  );
});
