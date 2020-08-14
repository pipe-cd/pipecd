import {
  CircularProgress,
  Link,
  makeStyles,
  Paper,
  Typography,
} from "@material-ui/core";
import CancelIcon from "@material-ui/icons/Cancel";
import dayjs from "dayjs";
import React, { FC, memo } from "react";
import { useDispatch, useSelector } from "react-redux";
import { Link as RouterLink } from "react-router-dom";
import { PAGE_PATH_APPLICATIONS } from "../constants";
import { DEPLOYMENT_STATE_TEXT } from "../constants/deployment-status-text";
import { AppState } from "../modules";
import { ActiveStage } from "../modules/active-stage";
import {
  cancelDeployment,
  Deployment,
  isDeploymentRunning,
  selectById as selectDeploymentById,
} from "../modules/deployments";
import {
  Environment,
  selectById as selectEnvById,
} from "../modules/environments";
import { Piped, selectById } from "../modules/pipeds";
import { fetchStageLog } from "../modules/stage-logs";
import { useInterval } from "../utils/use-interval";
import { StatusIcon } from "./deployment-status-icon";
import OpenInNewIcon from "@material-ui/icons/OpenInNew";
import { LabeledText } from "./labeled-text";
import { SplitButton } from "./split-button";

const useStyles = makeStyles((theme) => ({
  root: {
    padding: theme.spacing(2),
    position: "relative",
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
  deploymentMainInfo: {
    display: "flex",
    alignItems: "center",
  },
  contents: {
    display: "flex",
    flexDirection: "column",
  },
  detail: {
    display: "flex",
  },
  age: {
    color: theme.palette.text.secondary,
    marginLeft: theme.spacing(1),
  },
  content: {
    flex: 1,
  },
  commitInfo: {
    display: "flex",
  },
  actionButtons: {
    color: theme.palette.error.main,
    position: "absolute",
    top: theme.spacing(2),
    right: theme.spacing(2),
  },
  buttonProgress: {
    color: theme.palette.primary.main,
    position: "absolute",
    top: "50%",
    left: "50%",
    marginTop: -12,
    marginLeft: -12,
  },
  statusReason: {
    paddingTop: theme.spacing(1),
    paddingBottom: theme.spacing(1),
  },
  linkIcon: {
    fontSize: 16,
    verticalAlign: "text-bottom",
    marginLeft: theme.spacing(0.5),
  },
}));

interface Props {
  deploymentId: string;
}

const CANCEL_OPTIONS = [
  "Cancel",
  "Cancel with Rollback",
  "Cancel without Rollback",
];
const LOG_FETCH_INTERVAL = 2000;

export const DeploymentDetail: FC<Props> = memo(function DeploymentDetail({
  deploymentId,
}) {
  const classes = useStyles();
  const dispatch = useDispatch();

  const deployment = useSelector<AppState, Deployment | undefined>((state) =>
    selectDeploymentById(state.deployments, deploymentId)
  );
  const env = useSelector<AppState, Environment | undefined>((state) =>
    selectEnvById(state.environments, deployment?.envId || "")
  );
  const piped = useSelector<AppState, Piped | undefined>((state) => {
    return deployment
      ? selectById(state.pipeds, deployment.pipedId)
      : undefined;
  });
  const isCanceling = useSelector<AppState, boolean>((state) => {
    if (deployment?.id) {
      return state.deployments.canceling[deployment.id];
    }
    return false;
  });
  const activeStage = useSelector<AppState, ActiveStage | null>(
    (state) => state.activeStage
  );

  useInterval(
    () => {
      if (activeStage) {
        dispatch(
          fetchStageLog({
            deploymentId: activeStage.deploymentId,
            stageId: activeStage.stageId,
            offsetIndex: 0,
            retriedCount: 0,
          })
        );
      }
    },
    activeStage && isDeploymentRunning(deployment?.status)
      ? LOG_FETCH_INTERVAL
      : null
  );

  if (!deployment || !env || !piped) {
    return (
      <div className={classes.loading}>
        <CircularProgress />
      </div>
    );
  }

  return (
    <Paper square elevation={1} className={classes.root}>
      <div className={classes.contents}>
        <div className={classes.content}>
          <div className={classes.deploymentMainInfo}>
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
            <Typography variant="body1" className={classes.age}>
              {dayjs(deployment.createdAt * 1000).fromNow()}
            </Typography>
          </div>
          <Typography
            variant="body2"
            color="textSecondary"
            className={classes.statusReason}
          >
            {deployment.statusReason}
          </Typography>
        </div>
        <div className={classes.detail}>
          <div className={classes.content}>
            <LabeledText
              label="Application"
              value={
                <Link
                  variant="body2"
                  component={RouterLink}
                  to={`${PAGE_PATH_APPLICATIONS}/${deployment.applicationId}`}
                >
                  {deployment.applicationName}
                </Link>
              }
            />
            <LabeledText label="Piped" value={piped.name} />
            <LabeledText label="Summary" value={deployment.summary} />
          </div>
          <div className={classes.content}>
            {deployment.trigger.commit && (
              <div className={classes.commitInfo}>
                <Typography variant="subtitle2" color="textSecondary">
                  Commit:
                </Typography>
                <Typography variant="body2" className={classes.textMargin}>
                  {deployment.trigger.commit.message}
                </Typography>
                <span className={classes.textMargin}>
                  (
                  <Link
                    variant="body2"
                    href={deployment.trigger.commit.url}
                    target="_blank"
                    rel="noreferrer"
                  >
                    {`${deployment.trigger.commit.hash}`}
                    <OpenInNewIcon className={classes.linkIcon} />
                  </Link>
                  )
                </span>
              </div>
            )}
            <LabeledText
              label="Triggered by"
              value={
                deployment.trigger.commander ||
                deployment.trigger.commit?.author ||
                ""
              }
            />
          </div>
          {isDeploymentRunning(deployment.status) && (
            <SplitButton
              className={classes.actionButtons}
              options={CANCEL_OPTIONS}
              onClick={(index) => {
                dispatch(
                  cancelDeployment({
                    deploymentId,
                    forceRollback: index === 1,
                    forceNoRollback: index === 2,
                  })
                );
              }}
              startIcon={<CancelIcon />}
              loading={isCanceling}
            />
          )}
        </div>
      </div>
    </Paper>
  );
});
