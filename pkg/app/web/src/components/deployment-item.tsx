import React, { FC, memo } from "react";
import { makeStyles, Typography, ListItem } from "@material-ui/core";
import { useSelector } from "react-redux";
import { AppState } from "../modules";
import {
  Deployment,
  selectById as selectDeploymentById,
} from "../modules/deployments";
import { StatusIcon } from "./deployment-status-icon";
import {
  Application,
  selectById as selectApplicationById,
} from "../modules/applications";
import dayjs from "dayjs";
import { selectById, Environment } from "../modules/environments";
import { Link as RouterLink } from "react-router-dom";
import { PAGE_PATH_DEPLOYMENTS } from "../constants";

const useStyles = makeStyles((theme) => ({
  root: {
    flex: 1,
    padding: theme.spacing(2),
    display: "flex",
    alignItems: "center",
    height: 72,
    backgroundColor: theme.palette.background.paper,
  },
  env: {
    marginLeft: theme.spacing(1),
  },
  statusIcon: {
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
  const application = useSelector<AppState, Application | undefined>(
    (state) => {
      return deployment
        ? selectApplicationById(state.applications, deployment.applicationId)
        : undefined;
    }
  );
  const env = useSelector<AppState, Environment | undefined>((state) => {
    return deployment
      ? selectById(state.environments, deployment.envId)
      : undefined;
  });

  if (!deployment || !application || !env) {
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
          <Typography variant="h6">{application.name}</Typography>
          <Typography className={classes.env}>{env.name}</Typography>
          <StatusIcon
            status={deployment.status}
            className={classes.statusIcon}
          />
        </div>
        <Typography variant="body1" className={classes.description}>
          {deployment.summary}
        </Typography>
      </div>
      <div>{dayjs(deployment.createdAt * 1000).fromNow()}</div>
    </ListItem>
  );
});
