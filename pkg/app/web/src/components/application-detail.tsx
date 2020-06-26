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
import { Link as RouterLink } from "react-router-dom";
import { PAGE_PATH_DEPLOYMENTS } from "../constants";
import { APPLICATION_SYNC_STATUS_TEXT } from "../constants/application-sync-status-text";
import { AppState } from "../modules";
import {
  Application,
  selectById as selectApplicationById,
  ApplicationSyncStatus,
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
  textMargin: {
    marginLeft: theme.spacing(1),
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
    color: theme.palette.text.hint,
  },
}));

interface Props {
  applicationId: string;
}

export const ApplicationDetail: FC<Props> = memo(function ApplicationDetail({
  applicationId,
}) {
  const classes = useStyles();
  const app = useSelector<AppState, Application | undefined>((state) =>
    selectApplicationById(state.applications, applicationId)
  );
  const liveState = useSelector<AppState, ApplicationLiveState | undefined>(
    (state) => selectLiveStateById(state.applicationLiveState, applicationId)
  );
  const env = useSelector<AppState, Environment | undefined>((state) =>
    app ? selectEnvById(state.environments, app.envId) : undefined
  );

  if (!liveState || !app || !env) {
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
            <Typography variant="h5" className={classes.textMargin}>
              {app.name}
            </Typography>
            <Typography variant="subtitle2" className={classes.env}>
              {env.name}
            </Typography>
          </div>
          <div className={classes.statusLine}>
            <SyncStatusIcon status={app.syncState.status} />
            <div className={classes.statusText}>
              <Typography variant="h6" className={classes.syncStatusText}>
                {APPLICATION_SYNC_STATUS_TEXT[app.syncState.status]}
              </Typography>
              {app.syncState.status !== ApplicationSyncStatus.SYNCED && (
                <Typography
                  variant="body2"
                  className={classes.syncReason}
                >{`${app.syncState.shortReason}`}</Typography>
              )}
            </div>
          </div>
          <Typography className={classes.env} variant="body1">
            {dayjs(liveState.version.timestamp * 1000).fromNow()}
          </Typography>
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
      </Box>
    </Paper>
  );
});
