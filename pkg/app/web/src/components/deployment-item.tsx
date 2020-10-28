import { ListItem, makeStyles, Typography } from "@material-ui/core";
import dayjs from "dayjs";
import React, { FC, memo } from "react";
import { useSelector } from "react-redux";
import { Link as RouterLink } from "react-router-dom";
import { PAGE_PATH_DEPLOYMENTS } from "../constants/path";
import { DEPLOYMENT_STATE_TEXT } from "../constants/deployment-status-text";
import { AppState } from "../modules";
import {
  Deployment,
  selectById as selectDeploymentById,
} from "../modules/deployments";
import { Environment, selectById } from "../modules/environments";
import { StatusIcon } from "./deployment-status-icon";

const useStyles = makeStyles((theme) => ({
  root: {
    flex: 1,
    padding: theme.spacing(2),
    display: "flex",
    alignItems: "center",
    height: 72,
    backgroundColor: theme.palette.background.paper,
  },
  name: {
    marginLeft: theme.spacing(1),
  },
  env: {
    marginLeft: theme.spacing(1),
    color: theme.palette.text.secondary,
  },
  statusText: {
    marginLeft: theme.spacing(1),
  },
  commitHash: {
    marginLeft: theme.spacing(1),
  },
  head: {
    display: "flex",
    alignItems: "center",
  },
  description: {
    color: theme.palette.text.hint,
  },
  main: {
    flex: 1,
  },
}));

interface Props {
  id: string;
}

export const DeploymentItem: FC<Props> = memo(function DeploymentItem({ id }) {
  const classes = useStyles();
  const deployment = useSelector<AppState, Deployment | undefined>((state) =>
    selectDeploymentById(state.deployments, id)
  );

  const env = useSelector<AppState, Environment | undefined>((state) => {
    return deployment
      ? selectById(state.environments, deployment.envId)
      : undefined;
  });

  if (!deployment || !env) {
    return null;
  }

  return (
    <ListItem
      className={classes.root}
      button
      dense
      divider
      component={RouterLink}
      to={`${PAGE_PATH_DEPLOYMENTS}/${deployment.id}`}
    >
      <div className={classes.main}>
        <div className={classes.head}>
          <StatusIcon status={deployment.status} />
          <Typography variant="body1" className={classes.statusText}>
            {DEPLOYMENT_STATE_TEXT[deployment.status]}
          </Typography>
          <Typography variant="h6" className={classes.name}>
            {deployment.applicationName}
          </Typography>
          <Typography className={classes.env}>{env.name}</Typography>
        </div>
        <Typography variant="body1" className={classes.description}>
          {deployment.summary}
        </Typography>
      </div>
      <div>{dayjs(deployment.createdAt * 1000).fromNow()}</div>
    </ListItem>
  );
});
