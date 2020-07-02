import React, { FC, memo } from "react";
import { makeStyles } from "@material-ui/core";
import {
  HealthStatus,
  ApplicationLiveStateSnapshotModel,
} from "../modules/applications-live-state";
import FavoriteIcon from "@material-ui/icons/Favorite";
import OtherIcon from "@material-ui/icons/HelpOutline";
import UnknownIcon from "@material-ui/icons/ErrorOutline";

const useStyles = makeStyles((theme) => ({
  healthy: {
    color: theme.palette.success.main,
  },
  unknown: {
    color: theme.palette.warning.main,
  },
  other: {
    color: theme.palette.info.main,
  },
}));

export const KubernetesResourceHealthStatusIcon: FC<{
  health: HealthStatus;
}> = memo(function HealthStatusIcon({ health }) {
  const classes = useStyles();
  switch (health) {
    case HealthStatus.UNKNOWN:
      return <UnknownIcon fontSize="small" className={classes.unknown} />;
    case HealthStatus.HEALTHY:
      return <FavoriteIcon fontSize="small" className={classes.healthy} />;
    case HealthStatus.OTHER:
      return <OtherIcon fontSize="small" className={classes.other} />;
  }
});

export const ApplicationHealthStatusIcon: FC<{
  health: ApplicationLiveStateSnapshotModel.Status;
}> = memo(function HealthStatusIcon({ health }) {
  const classes = useStyles();
  switch (health) {
    case ApplicationLiveStateSnapshotModel.Status.UNKNOWN:
      return <UnknownIcon className={classes.unknown} />;
    case ApplicationLiveStateSnapshotModel.Status.HEALTHY:
      return <FavoriteIcon className={classes.healthy} />;
    case ApplicationLiveStateSnapshotModel.Status.OTHER:
      return <OtherIcon className={classes.other} />;
  }
});
