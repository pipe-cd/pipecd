import React, { FC } from "react";
import { makeStyles, Paper, Typography, Box } from "@material-ui/core";
import { LabeledText } from "./labeled-text";
import { ApplicationSyncStatus } from "pipe/pkg/app/web/model/application_pb";
import { SyncStatusIcon } from "./sync-status-icon";
import { APPLICATION_SYNC_STATUS_TEXT } from "../constants/application-sync-status-text";

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
  name: string;
  env: string;
  piped: string;
  status: ApplicationSyncStatus;
  deployedAt: number;
  version: string;
  deploymentId: string;
  description: string;
}

export const ApplicationDetail: FC<Props> = ({
  name,
  env,
  piped,
  status,
  deployedAt,
  description,
  deploymentId,
  version,
}) => {
  const classes = useStyles();
  return (
    <Paper square elevation={1} className={classes.container}>
      <Box display="flex">
        <Box flex={1}>
          <Box display="flex" alignItems="center">
            <Typography variant="h6" className={classes.textMargin}>
              {name}
            </Typography>
            <Typography variant="subtitle2" className={classes.env}>
              {env}
            </Typography>
          </Box>
          <Box display="flex">
            <SyncStatusIcon status={status} />
            <Typography variant="subtitle1" className={classes.syncStatusText}>
              {APPLICATION_SYNC_STATUS_TEXT[status]}
            </Typography>

            <SyncStatusIcon status={status} />
            <Typography variant="subtitle1" className={classes.statusText}>
              {/** TODO: Show health status */}
              Healthy
            </Typography>
          </Box>
          <Typography className={classes.env} variant="body1">
            {deployedAt}
          </Typography>
        </Box>
        <Box flex={1}>
          <LabeledText label="Latest Deployment" text={deploymentId} />
          <LabeledText label="Version" text={version} />
          <LabeledText label="Description" text={description} />
        </Box>
      </Box>
    </Paper>
  );
};
