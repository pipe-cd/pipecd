import { makeStyles } from "@material-ui/core";
import { Cached, CheckCircle, Error, Info, ReportOff } from "@material-ui/icons";
import { FC } from "react";
import { ApplicationSyncStatus } from "~/modules/applications";

const useStyles = makeStyles((theme) => ({
  [ApplicationSyncStatus.UNKNOWN]: {
    color: theme.palette.grey[500],
  },
  [ApplicationSyncStatus.SYNCED]: {
    color: theme.palette.success.main,
  },
  [ApplicationSyncStatus.DEPLOYING]: {
    color: theme.palette.info.main,
    animation: `$running 3s linear infinite`,
  },
  [ApplicationSyncStatus.OUT_OF_SYNC]: {
    color: theme.palette.error.main,
  },
  [ApplicationSyncStatus.INVALID_CONFIG]: {
    color: theme.palette.error.light,
  },
  "@keyframes running": {
    "0%": {
      transform: "rotate(360deg)",
    },
    "100%": {
      transform: "rotate(0deg)",
    },
  },
}));

export interface SyncStatusIconProps {
  status: ApplicationSyncStatus;
}

export const SyncStatusIcon: FC<SyncStatusIconProps> = ({ status }) => {
  const classes = useStyles();

  switch (status) {
    case ApplicationSyncStatus.UNKNOWN:
      return <Info className={classes[status]} />;
    case ApplicationSyncStatus.SYNCED:
      return <CheckCircle className={classes[status]} />;
    case ApplicationSyncStatus.DEPLOYING:
      return <Cached className={classes[status]} />;
    case ApplicationSyncStatus.OUT_OF_SYNC:
      return <Error className={classes[status]} />;
    case ApplicationSyncStatus.INVALID_CONFIG:
      return <ReportOff className={classes[status]} />;
  }
};
