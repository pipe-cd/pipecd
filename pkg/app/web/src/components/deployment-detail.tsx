import {
  Box,
  Link,
  makeStyles,
  Paper,
  Typography,
  CircularProgress,
} from "@material-ui/core";
import dayjs from "dayjs";
import React, { FC, memo } from "react";
import { useSelector } from "react-redux";
import {
  selectById as selectDeploymentById,
  Deployment,
} from "../modules/deployments";
import {
  Environment,
  selectById as selectEnvById,
} from "../modules/environments";
import { DEPLOYMENT_STATE_TEXT } from "../constants/deployment-status-text";
import { AppState } from "../modules";
import { StatusIcon } from "./deployment-status-icon";
import { LabeledText } from "./labeled-text";

const useStyles = makeStyles((theme) => ({
  container: {
    padding: theme.spacing(3),
  },
  textMargin: {
    marginLeft: theme.spacing(1),
  },
  env: {
    color: theme.palette.text.secondary,
    marginLeft: theme.spacing(1),
  },
  loading: {
    flex: 1,
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
  },
}));

interface Props {
  deploymentId: string;
}

export const DeploymentDetail: FC<Props> = memo(function DeploymentDetail({
  deploymentId,
}) {
  const classes = useStyles();

  const deployment = useSelector<AppState, Deployment | undefined>((state) =>
    selectDeploymentById(state.deployments, deploymentId)
  );
  const env = useSelector<AppState, Environment | undefined>((state) =>
    selectEnvById(state.environments, deployment?.envId || "")
  );

  if (!deployment || !env) {
    return (
      <div className={classes.loading}>
        <CircularProgress />
      </div>
    );
  }

  return (
    <Paper square elevation={1} className={classes.container}>
      <Box display="flex">
        <Box flex={1}>
          <Box alignItems="center" display="flex">
            <StatusIcon status={deployment.status} />
            <Typography className={classes.textMargin} variant="h6">
              {DEPLOYMENT_STATE_TEXT[deployment.status]}
            </Typography>
            <Typography className={classes.textMargin} variant="h6">
              {deployment.id}
            </Typography>
            <Typography variant="subtitle1" className={classes.env}>
              {env.name}
            </Typography>
          </Box>
          <Typography variant="subtitle1">
            {dayjs(deployment.createdAt * 1000).fromNow()}
          </Typography>

          <LabeledText label="Piped ID" value={deployment.pipedId} />
          <LabeledText label="Description" value={deployment.description} />
        </Box>
        <Box flex={1}>
          {deployment.trigger.commit && (
            <Box display="flex">
              <Typography variant="subtitle2" color="textSecondary">
                COMMIT:
              </Typography>
              <Box display="flex">
                <Typography variant="body2" className={classes.textMargin}>
                  {deployment.trigger.commit.message}
                </Typography>
                <span className={classes.textMargin}>
                  (
                  <Link variant="body2">{`${deployment.trigger.commit.hash}`}</Link>
                  )
                </span>
              </Box>
            </Box>
          )}
          <LabeledText
            label="TRIGGERED BY"
            value={deployment.trigger.commander}
          />
        </Box>
      </Box>
    </Paper>
  );
});
