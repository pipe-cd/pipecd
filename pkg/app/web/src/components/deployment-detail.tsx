import {
  Link,
  makeStyles,
  Paper,
  Typography,
  CircularProgress,
  Button,
} from "@material-ui/core";
import dayjs from "dayjs";
import React, { FC, memo } from "react";
import { useSelector, useDispatch } from "react-redux";
import {
  selectById as selectDeploymentById,
  Deployment,
  cancelDeployment,
  isDeploymentRunning,
} from "../modules/deployments";
import {
  Environment,
  selectById as selectEnvById,
} from "../modules/environments";
import { DEPLOYMENT_STATE_TEXT } from "../constants/deployment-status-text";
import { AppState } from "../modules";
import { StatusIcon } from "./deployment-status-icon";
import { LabeledText } from "./labeled-text";
import CancelIcon from "@material-ui/icons/Cancel";

const useStyles = makeStyles((theme) => ({
  root: {
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
  deploymentMainInfo: {
    display: "flex",
    alignItems: "center",
  },
  contents: {
    display: "flex",
  },
  content: {
    flex: 1,
  },
  commitInfo: {
    display: "flex",
  },
  buttonArea: {
    color: theme.palette.error.main,
  },
  buttonProgress: {
    color: theme.palette.primary.main,
    position: "absolute",
    top: "50%",
    left: "50%",
    marginTop: -12,
    marginLeft: -12,
  },
}));

interface Props {
  deploymentId: string;
}

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
  const isCanceling = useSelector<AppState, boolean>((state) => {
    if (deployment?.id) {
      return state.deployments.canceling[deployment.id];
    }
    return false;
  });

  const handleCancel = (): void => {
    dispatch(cancelDeployment({ deploymentId, withoutRollback: false }));
  };

  if (!deployment || !env) {
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
          </div>
          <Typography variant="subtitle1">
            {dayjs(deployment.createdAt * 1000).fromNow()}
          </Typography>

          <LabeledText label="Piped ID" value={deployment.pipedId} />
          <LabeledText label="Description" value={deployment.description} />
        </div>
        <div className={classes.content}>
          {deployment.trigger.commit && (
            <div className={classes.commitInfo}>
              <Typography variant="subtitle2" color="textSecondary">
                COMMIT:
              </Typography>
              <Typography variant="body2" className={classes.textMargin}>
                {deployment.trigger.commit.message}
              </Typography>
              <span className={classes.textMargin}>
                (
                <Link variant="body2">{`${deployment.trigger.commit.hash}`}</Link>
                )
              </span>
            </div>
          )}
          <LabeledText
            label="TRIGGERED BY"
            value={deployment.trigger.commander}
          />
        </div>
        {isDeploymentRunning(deployment.status) && (
          <div className={classes.buttonArea}>
            <Button
              color="inherit"
              variant="outlined"
              startIcon={<CancelIcon />}
              onClick={handleCancel}
              disabled={isCanceling}
            >
              CANCEL
              {isCanceling && (
                <CircularProgress
                  size={24}
                  className={classes.buttonProgress}
                />
              )}
            </Button>
          </div>
        )}
      </div>
    </Paper>
  );
});
