import { Box, Typography } from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";
import { FC } from "react";
import { APPLICATION_SYNC_STATUS_TEXT } from "~/constants/application-sync-status-text";
import { UI_TEXT_NOT_AVAILABLE_TEXT } from "~/constants/ui-text";
import {
  ApplicationSyncState,
  ApplicationSyncStatus,
} from "~/modules/applications";
import { SyncStatusIcon } from "./sync-status-icon";

const useStyles = makeStyles((theme) => ({
  statusText: {
    marginLeft: theme.spacing(0.5),
    whiteSpace: "nowrap",
  },
}));

export interface AppSyncStatusProps {
  syncState?: ApplicationSyncState.AsObject;
  deploying: boolean;
  className?: string;
  size?: "medium" | "large";
}

export const AppSyncStatus: FC<AppSyncStatusProps> = ({
  syncState,
  deploying,
  className,
  size = "medium",
}) => {
  const classes = useStyles();
  const fontVariant = size === "medium" ? "body2" : "h6";

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
