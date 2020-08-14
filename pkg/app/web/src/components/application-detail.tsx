import {
  Button,
  CircularProgress,
  Link,
  makeStyles,
  Paper,
  Typography,
} from "@material-ui/core";
import SyncIcon from "@material-ui/icons/Cached";
import OpenInNewIcon from "@material-ui/icons/OpenInNew";
import dayjs from "dayjs";
import React, { FC, memo } from "react";
import { useDispatch, useSelector } from "react-redux";
import { Link as RouterLink } from "react-router-dom";
import { PAGE_PATH_DEPLOYMENTS } from "../constants";
import { APPLICATION_SYNC_STATUS_TEXT } from "../constants/application-sync-status-text";
import { APPLICATION_HEALTH_STATUS_TEXT } from "../constants/health-status-text";
import { AppState } from "../modules";
import {
  Application,
  selectById as selectApplicationById,
  syncApplication,
} from "../modules/applications";
import {
  ApplicationLiveState,
  selectById as selectLiveStateById,
} from "../modules/applications-live-state";
import {
  Environment,
  selectById as selectEnvById,
} from "../modules/environments";
import { Piped, selectById as selectPipeById } from "../modules/pipeds";
import { ApplicationHealthStatusIcon } from "./health-status-icon";
import { LabeledText } from "./labeled-text";
import { SyncStateReason } from "./sync-state-reason";
import { SyncStatusIcon } from "./sync-status-icon";

const useStyles = makeStyles((theme) => ({
  root: {
    padding: theme.spacing(2),
    display: "flex",
    zIndex: theme.zIndex.appBar,
    position: "relative",
    flexDirection: "column",
  },
  nameAndEnv: {
    display: "flex",
    alignItems: "baseline",
  },
  mainContent: { flex: 1 },
  content: {
    flex: 1,
  },
  detail: {
    display: "flex",
    marginTop: theme.spacing(1),
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
  actionButtons: {
    position: "absolute",
    right: theme.spacing(2),
    top: theme.spacing(2),
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
  linkIcon: {
    fontSize: 16,
    verticalAlign: "text-bottom",
    marginLeft: theme.spacing(0.5),
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
  const pipe = useSelector<AppState, Piped | undefined>((state) =>
    app ? selectPipeById(state.pipeds, app.pipedId) : undefined
  );

  const handleSync = (): void => {
    if (app) {
      dispatch(syncApplication({ applicationId: app.id }));
    }
  };

  if (!app || !env || !pipe) {
    return (
      <div className={classes.loading}>
        <CircularProgress />
      </div>
    );
  }

  return (
    <Paper square elevation={1} className={classes.root}>
      <div className={classes.mainContent}>
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

        {app.syncState && (
          <>
            <div className={classes.statusLine}>
              <SyncStatusIcon status={app.syncState.status} />
              <div className={classes.statusText}>
                <Typography variant="h6" className={classes.syncStatusText}>
                  {APPLICATION_SYNC_STATUS_TEXT[app.syncState.status]}
                </Typography>
              </div>

              {liveState && (
                <>
                  <ApplicationHealthStatusIcon
                    health={liveState.healthStatus}
                  />
                  <Typography variant="h6" className={classes.syncStatusText}>
                    {APPLICATION_HEALTH_STATUS_TEXT[liveState.healthStatus]}
                  </Typography>
                </>
              )}
            </div>

            <SyncStateReason
              summary={app.syncState.shortReason}
              detail={app.syncState.reason}
            />
          </>
        )}
      </div>

      <div className={classes.detail}>
        <div className={classes.content}>
          <LabeledText label="piped" value={`${pipe.name}`} />

          <LabeledText label="Cloud Provider" value={`${app.cloudProvider}`} />

          <LabeledText
            label="Git Path"
            value={
              <Link href={app.gitPath.url} target="_blank" rel="noreferrer">
                {`${app.gitPath.repo !== undefined ? app.gitPath.repo.id : ""}`}
                <OpenInNewIcon className={classes.linkIcon} />
              </Link>
            }
          />
        </div>

        {app.mostRecentlySuccessfulDeployment && (
          <div className={classes.content}>
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
          </div>
        )}
      </div>

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
    </Paper>
  );
});
