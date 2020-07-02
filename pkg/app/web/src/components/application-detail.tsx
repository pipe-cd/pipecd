import {
  Box,
  Link,
  makeStyles,
  Paper,
  Typography,
  CircularProgress,
  Button,
} from "@material-ui/core";
import dayjs from "dayjs";
import React, { FC, memo, useState } from "react";
import { useSelector, useDispatch } from "react-redux";
import { Link as RouterLink } from "react-router-dom";
import { PAGE_PATH_DEPLOYMENTS } from "../constants";
import { APPLICATION_SYNC_STATUS_TEXT } from "../constants/application-sync-status-text";
import { AppState } from "../modules";
import {
  Application,
  selectById as selectApplicationById,
  ApplicationSyncStatus,
  syncApplication,
} from "../modules/applications";
import {
  ApplicationLiveState,
  selectById as selectLiveStateById,
} from "../modules/applications-live-state";
import { LabeledText } from "./labeled-text";
import { SyncStatusIcon } from "./sync-status-icon";
import {
  selectById as selectEnvById,
  Environment,
} from "../modules/environments";
import SyncIcon from "@material-ui/icons/Cached";
import { ApplicationHealthStatusIcon } from "./health-status-icon";
import { APPLICATION_HEALTH_STATUS_TEXT } from "../constants/health-status-text";

const useStyles = makeStyles((theme) => ({
  nameAndEnv: {
    display: "flex",
    alignItems: "baseline",
  },
  container: {
    padding: theme.spacing(2),
  },
  loading: {
    flex: 1,
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
  },

  env: {
    color: theme.palette.text.secondary,
    marginLeft: theme.spacing(1),
  },
  statusLine: {
    display: "flex",
    alignItems: "center",
  },
  statusText: {
    display: "flex",
    alignItems: "baseline",
  },
  syncStatusText: {
    marginLeft: theme.spacing(0.5),
    marginRight: theme.spacing(1),
  },
  syncReason: {
    color: theme.palette.text.secondary,
  },
  learnMore: {
    color: theme.palette.primary.light,
    marginLeft: theme.spacing(1),
  },
  reasonDetail: {
    padding: theme.spacing(2),
    fontFamily: "Roboto Mono",
    marginTop: theme.spacing(1),
    wordBreak: "break-all",
  },
  reasonLine: {
    display: "flex",
    alignItems: "center",
  },
  actionButtons: {
    paddingLeft: theme.spacing(2),
    paddingRight: theme.spacing(2),
  },
  buttonProgress: {
    color: theme.palette.primary.main,
    position: "absolute",
    top: "50%",
    left: "50%",
    marginTop: -12,
    marginLeft: -12,
  },
  age: {
    color: theme.palette.text.secondary,
    marginLeft: theme.spacing(1),
  },
}));

interface Props {
  applicationId: string;
}

const useIsSyncingApplication = (applicationId: string | null): boolean => {
  return useSelector<AppState, boolean>((state) => {
    if (!applicationId) {
      return false;
    }

    return state.applications.syncing[applicationId];
  });
};

export const ApplicationDetail: FC<Props> = memo(function ApplicationDetail({
  applicationId,
}) {
  const classes = useStyles();
  const [showReason, setShowReason] = useState(false);
  const dispatch = useDispatch();
  const app = useSelector<AppState, Application | undefined>((state) =>
    selectApplicationById(state.applications, applicationId)
  );
  const liveState = useSelector<AppState, ApplicationLiveState | undefined>(
    (state) => selectLiveStateById(state.applicationLiveState, applicationId)
  );
  const env = useSelector<AppState, Environment | undefined>((state) =>
    app ? selectEnvById(state.environments, app.envId) : undefined
  );
  const isSyncing = useIsSyncingApplication(app ? app.id : null);

  const handleSync = (): void => {
    if (app) {
      dispatch(syncApplication({ applicationId: app.id }));
    }
  };

  if (!app || !env) {
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
          <div className={classes.nameAndEnv}>
            <Typography variant="h5">{app.name}</Typography>
            <Typography variant="subtitle2" className={classes.env}>
              {env.name}
            </Typography>

            {liveState && (
              <Typography className={classes.age} variant="body1">
                {dayjs(liveState.version.timestamp * 1000).fromNow()}
              </Typography>
            )}
          </div>

          <div className={classes.statusLine}>
            <SyncStatusIcon status={app.syncState.status} />
            <div className={classes.statusText}>
              <Typography variant="h6" className={classes.syncStatusText}>
                {APPLICATION_SYNC_STATUS_TEXT[app.syncState.status]}
              </Typography>
            </div>

            {liveState && (
              <>
                <ApplicationHealthStatusIcon health={liveState.healthStatus} />
                <Typography variant="h6" className={classes.syncStatusText}>
                  {APPLICATION_HEALTH_STATUS_TEXT[liveState.healthStatus]}
                </Typography>
              </>
            )}
          </div>

          {app.syncState.status !== ApplicationSyncStatus.SYNCED && (
            <div className={classes.reasonLine}>
              <Typography variant="body2" className={classes.syncReason}>
                {app.syncState.shortReason}
              </Typography>
              {app.syncState.shortReason && (
                <Button
                  variant="text"
                  size="small"
                  className={classes.learnMore}
                  onClick={() => setShowReason(!showReason)}
                >
                  {showReason ? "HIDE DETAIL" : "SHOW DETAIL"}
                </Button>
              )}
            </div>
          )}

          {showReason && (
            <Paper
              elevation={0}
              variant="outlined"
              className={classes.reasonDetail}
            >
              {app.syncState.reason.split("\n").map((line, i) => (
                <div key={i}>{line}</div>
              ))}
            </Paper>
          )}
        </Box>
        <Box flex={1}>
          <LabeledText
            label="Latest Deployment"
            value={
              <Link
                component={RouterLink}
                to={`${PAGE_PATH_DEPLOYMENTS}/${app.mostRecentlySuccessfulDeployment.deploymentId}`}
              >
                {app.mostRecentlySuccessfulDeployment.deploymentId}
              </Link>
            }
          />
          <LabeledText
            label="Version"
            value={app.mostRecentlySuccessfulDeployment.version}
          />
          <LabeledText
            label="Description"
            value={app.mostRecentlySuccessfulDeployment.description}
          />
        </Box>
        <div className={classes.actionButtons}>
          <Button
            variant="outlined"
            color="primary"
            onClick={handleSync}
            disabled={isSyncing}
            startIcon={<SyncIcon />}
          >
            SYNC
            {isSyncing && (
              <CircularProgress size={24} className={classes.buttonProgress} />
            )}
          </Button>
        </div>
      </Box>
    </Paper>
  );
});
