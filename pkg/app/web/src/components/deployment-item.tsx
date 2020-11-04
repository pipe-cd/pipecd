import { Box, ListItem, makeStyles, Typography } from "@material-ui/core";
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
import { APPLICATION_KIND_TEXT } from "../constants/application-kind";
import { ellipsis } from "../styles/text";

const useStyles = makeStyles((theme) => ({
  root: {
    flex: 1,
    padding: theme.spacing(2),
    display: "flex",
    alignItems: "center",
    height: 72,
    backgroundColor: theme.palette.background.paper,
  },
  info: {
    marginLeft: theme.spacing(1),
  },
  statusText: {
    marginLeft: theme.spacing(1),
    lineHeight: "1.5rem",
    // Fix width to prevent misalignment of application name.
    width: "100px",
  },
  description: {
    ...ellipsis,
    color: theme.palette.text.hint,
  },
}));

interface Props {
  id: string;
}

const NO_DESCRIPTION = "No description.";

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
      <Box display="flex" alignItems="center">
        <StatusIcon status={deployment.status} />
        <Typography
          variant="subtitle2"
          className={classes.statusText}
          component="span"
        >
          {DEPLOYMENT_STATE_TEXT[deployment.status]}
        </Typography>
      </Box>
      <Box
        display="flex"
        flexDirection="column"
        flex={1}
        pl={2}
        overflow="hidden"
      >
        <Box display="flex" alignItems="baseline">
          <Typography variant="h6" component="span">
            {deployment.applicationName}
          </Typography>
          <Typography
            variant="subtitle2"
            className={classes.info}
            component="span"
          >
            {env.name}
          </Typography>
          <Typography
            variant="body2"
            color="textSecondary"
            className={classes.info}
          >
            {APPLICATION_KIND_TEXT[deployment.kind]}
          </Typography>
        </Box>
        <Typography variant="body1" className={classes.description}>
          {deployment.summary || NO_DESCRIPTION}
        </Typography>
      </Box>
      <div>{dayjs(deployment.createdAt * 1000).fromNow()}</div>
    </ListItem>
  );
});
