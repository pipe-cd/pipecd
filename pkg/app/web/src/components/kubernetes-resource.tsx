import React, { FC } from "react";
import { makeStyles, Paper, Typography } from "@material-ui/core";

const useStyles = makeStyles((theme) => ({
  container: {
    display: "inline-flex",
    flexDirection: "column",
    padding: theme.spacing(2),
  },
}));

interface Props {
  name: string;
  kind: string;
}

export const KubernetesResource: FC<Props> = ({ name, kind }) => {
  const classes = useStyles();
  return (
    <Paper square className={classes.container}>
      <Typography variant="caption">{kind}</Typography>
      <Typography variant="subtitle2">{name}</Typography>
    </Paper>
  );
};
