import React, { FC } from "react";
import { makeStyles } from "@material-ui/core";

const useStyles = makeStyles(() => ({
  container: {},
}));

interface Props {}

export const KubernetesStateView: FC<Props> = ({}) => {
  const classes = useStyles();
  return <div className={classes.container}>Hello</div>;
};
