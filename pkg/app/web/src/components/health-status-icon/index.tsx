import { FC, memo } from "react";
import { makeStyles } from "@material-ui/core";
import {
  HealthStatus,
  ApplicationLiveStateSnapshot,
} from "~/modules/applications-live-state";
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

export interface KubernetesResourceHealthStatusIconProps {
  health: HealthStatus;
}

export const KubernetesResourceHealthStatusIcon: FC<KubernetesResourceHealthStatusIconProps> = memo(
  function HealthStatusIcon({ health }) {
    const classes = useStyles();
    switch (health) {
      case HealthStatus.UNKNOWN:
        return <UnknownIcon fontSize="small" className={classes.unknown} />;
      case HealthStatus.HEALTHY:
        return <FavoriteIcon fontSize="small" className={classes.healthy} />;
      case HealthStatus.OTHER:
        return <OtherIcon fontSize="small" className={classes.other} />;
    }
  }
);

export interface ApplicationHealthStatusIconProps {
  health: ApplicationLiveStateSnapshot.Status;
}

export const ApplicationHealthStatusIcon: FC<ApplicationHealthStatusIconProps> = memo(
  function HealthStatusIcon({ health }) {
    const classes = useStyles();
    switch (health) {
      case ApplicationLiveStateSnapshot.Status.UNKNOWN:
        return <UnknownIcon className={classes.unknown} />;
      case ApplicationLiveStateSnapshot.Status.HEALTHY:
        return <FavoriteIcon className={classes.healthy} />;
      case ApplicationLiveStateSnapshot.Status.OTHER:
        return <OtherIcon className={classes.other} />;
    }
  }
);
