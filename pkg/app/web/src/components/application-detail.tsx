import { Box, Link, makeStyles, Paper, Typography } from "@material-ui/core";
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
} from "../modules/applications";
import {
  ApplicationLiveState,
  selectById as selectLiveStateById,
} from "../modules/applications-live-state";
import { LabeledText } from "./labeled-text";
import { SyncStatusIcon } from "./sync-status-icon";

const useStyles = makeStyles((theme) => ({
  container: {
    padding: theme.spacing(2),
  },
  textMargin: {
    marginLeft: theme.spacing(1),
  },
  syncStatusText: {
    marginRight: theme.spacing(1),
    marginLeft: theme.spacing(0.5),
  },
  statusText: {
    marginLeft: theme.spacing(0.5),
  },
  env: {
    color: theme.palette.text.secondary,
    marginLeft: theme.spacing(1),
  },
}));

interface Props {
  applicationId: string;
}

export const ApplicationDetail: FC<Props> = memo(({ applicationId }) => {
  const classes = useStyles();
  const app = useSelector<AppState, Application | undefined>((state) =>
    selectApplicationById(state.applications, applicationId)
  );
  const liveState = useSelector<AppState, ApplicationLiveState | undefined>(
    (state) => selectLiveStateById(state.applicationLiveState, applicationId)
  );

  if (!liveState || !app) {
    return null;
  }

  return (
    <Paper square elevation={1} className={classes.container}>
      <Box display="flex">
        <Box flex={1}>
          <Box display="flex" alignItems="center">
            <Typography variant="h6" className={classes.textMargin}>
              {app.name}
            </Typography>
            <Typography variant="subtitle2" className={classes.env}>
              {liveState.envId}
            </Typography>
          </Box>
          <Box display="flex">
            <SyncStatusIcon status={app.syncState.status} />
            <Typography variant="subtitle1" className={classes.syncStatusText}>
              {APPLICATION_SYNC_STATUS_TEXT[app.syncState.status]}
            </Typography>

            <SyncStatusIcon status={app.syncState.status} />
            <Typography variant="subtitle1" className={classes.statusText}>
              {/** TODO: Show health status */}
              Healthy
            </Typography>
          </Box>
          <Typography className={classes.env} variant="body1">
            {dayjs(liveState.version.timestamp).fromNow()}
          </Typography>
        </Box>
        <Box flex={1}>
          <LabeledText
            label="Latest Deployment"
            value={
              <Link
                component={RouterLink}
                to={`${PAGE_PATH_DEPLOYMENTS}/${app.mostRecentSuccessfulDeployment.deploymentId}`}
              >
                {app.mostRecentSuccessfulDeployment.deploymentId}
              </Link>
            }
          />
          <LabeledText label="Version" value={`${liveState.version.index}`} />
          <LabeledText label="Description" value="description" />
        </Box>
      </Box>
    </Paper>
  );
});
