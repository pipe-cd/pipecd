import React, { FC } from "react";
import { Box, makeStyles, Typography } from "@material-ui/core";
import { SyncStatusIcon } from "./sync-status-icon";
import { APPLICATION_SYNC_STATUS_TEXT } from "../constants/application-sync-status-text";
import {
  ApplicationSyncState,
  ApplicationSyncStatus,
} from "../modules/applications";
import { UI_TEXT_NOT_AVAILABLE_TEXT } from "../constants/ui-text";

const useStyles = makeStyles((theme) => ({
  statusText: {
    marginLeft: theme.spacing(0.5),
  },
}));

interface Props {
  syncState?: ApplicationSyncState;
  deploying: boolean;
  className?: string;
  size?: "medium" | "large";
}

export const AppSyncStatus: FC<Props> = ({
  syncState,
  deploying,
  className,
  size = "medium",
}) => {
  const classes = useStyles();
  const fontVariant = size === "medium" ? "body1" : "h6";

  if (deploying) {
    return (
      <Box display="flex" alignItems="center" className={className}>
        <SyncStatusIcon status={ApplicationSyncStatus.DEPLOYING} />
        <Typography
          className={classes.statusText}
          variant={fontVariant}
          component="span"
        >
          {APPLICATION_SYNC_STATUS_TEXT[ApplicationSyncStatus.DEPLOYING]}
        </Typography>
      </Box>
    );
  }

  return (
    <Box display="flex" alignItems="center" className={className}>
      {syncState ? <SyncStatusIcon status={syncState.status} /> : null}
      <Typography
        className={classes.statusText}
        variant={fontVariant}
        component="span"
      >
        {syncState
          ? APPLICATION_SYNC_STATUS_TEXT[syncState.status]
          : UI_TEXT_NOT_AVAILABLE_TEXT}
      </Typography>
    </Box>
  );
};
