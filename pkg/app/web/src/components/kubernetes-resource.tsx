import React, { FC } from "react";
import { makeStyles, Paper, Typography } from "@material-ui/core";
import { HealthStatus } from "../modules/applications-live-state";
import { HealthStatusIcon } from "./health-status-icon";

const useStyles = makeStyles((theme) => ({
  container: {
    display: "inline-flex",
    flexDirection: "column",
    padding: theme.spacing(2),
    width: 300,
  },
  nameLine: {
    display: "flex",
  },
  name: {
    marginLeft: theme.spacing(0.5),
  },
}));

interface Props {
  name: string;
  kind: string;
  health: HealthStatus;
}

export const KubernetesResource: FC<Props> = ({ name, kind, health }) => {
  const classes = useStyles();
  return (
    <Paper square className={classes.container}>
      <Typography variant="caption">{kind}</Typography>
      <div className={classes.nameLine}>
        <HealthStatusIcon health={health} />
        <Typography variant="subtitle2" className={classes.name}>
          {name}
        </Typography>
      </div>
    </Paper>
  );
};
