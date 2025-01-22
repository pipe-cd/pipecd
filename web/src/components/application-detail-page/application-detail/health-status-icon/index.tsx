import { makeStyles } from "@material-ui/core";
import UnknownIcon from "@material-ui/icons/ErrorOutline";
import FavoriteIcon from "@material-ui/icons/Favorite";
import OtherIcon from "@material-ui/icons/HelpOutline";
import { FC, memo } from "react";
import { ApplicationLiveStateSnapshot } from "~/modules/applications-live-state";

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
  unhealthy: {
    color: theme.palette.info.main,
  },
}));

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
      case ApplicationLiveStateSnapshot.Status.UNHEALTHY:
        return <OtherIcon className={classes.unhealthy} />;
    }
  }
);
